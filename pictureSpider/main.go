package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/joho/godotenv"
)

var targetURL = "https://www.jyflb.com/tag/%e9%82%aa%e6%81%b6%e5%8a%a8%e6%80%81%e5%9b%be/page/"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var startPage = flag.Int("start", 1, "page start")
	var count = flag.Int("length", 100, "download length")
	var resetJson = flag.Int("json", 0, "reset data.json default not reset")
	flag.Parse()
	if *resetJson == 1 {
		GenerateJson()
	} else {
		for i := *startPage; i < *count; i++ {
			haveNext := make(chan bool)
			nowURL := targetURL + strconv.Itoa(i)
			fmt.Println(nowURL)
			go fetchPage(nowURL, haveNext)
			if !<-haveNext {
				fmt.Println("this page has no more content to fetch")
				os.Exit(1)
			}
		}
	}
}

type GenJson struct {
	Title string   `json:"title"`
	Files []string `json:"files"`
}

func GenerateJson() {
	file, err := ioutil.ReadDir("images")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var jsons []GenJson
	for _, v := range file {
		var oneJson GenJson
		if v.IsDir() && v.Name() != "" {
			oneJson.Title = v.Name()
			image, err := ioutil.ReadDir("images/" + v.Name())
			if err != nil {
				fmt.Println(err)
			} else {
				for _, j := range image {
					if !j.IsDir() && j.Name() != "" {
						oneJson.Files = append(oneJson.Files, j.Name())
					}
				}
				if len(oneJson.Files) == 0 {
					fmt.Println("images/" + v.Name())
				}
			}
		}
		jsons = append(jsons, oneJson)
	}
	writeData, err := json.Marshal(jsons)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("data.json", writeData, 0777)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func filterDislike(spString string) bool {
	dislikeKeyWord := os.Getenv("dislike")
	dislikeKeyWords := strings.Split(dislikeKeyWord, ",")

	favouriteKeyword := os.Getenv("favourite")
	favouriteKeywords := strings.Split(favouriteKeyword, ",")
	//favourite priority is big
	if len(favouriteKeywords) >= 1 {
		for _, v := range favouriteKeywords {
			res := strings.Split(spString, v)
			if len(res) >= 2 {
				return false
			}
		}
		return true
	}
	for _, v := range dislikeKeyWords {
		res := strings.Split(spString, v)
		if len(res) >= 2 {
			return true
		}
	}
	return false
}

func fetchNextURL(nowURL string, chanNextURL chan string) {
	doc, err := goquery.NewDocument(nowURL)
	if err != nil {
		panic(err)
	}
	section := doc.Find(".content .pagination .next-page a")
	nextURL, _ := section.Attr("href")
	fmt.Println("next url = " + nextURL + "\n")
	chanNextURL <- nextURL
}

func fetchPage(url string, hasNext chan bool) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	count := make(chan int)
	ct := 0
	doc.Find(".excerpt").EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		title := contentSelection.Find("h2 a").Text()
		targetURL, _ := contentSelection.Find("h2 a").Attr("href")
		folderName := strings.Replace(title, " ", "", -1)
		isDislike := filterDislike(folderName)
		if !isDislike {
			folder := "./images/" + folderName
			createFolder, err := createFolder(folder)
			if err != nil {
				panic(err)
			}
			if createFolder {
				go fetchDetail(targetURL, folder, count)
				ct++
				allFinish := <-count
				fmt.Printf("title : %s download %d images \n", folderName, allFinish)
			} else {
				fmt.Printf("skip %s because exists\n", folder)
			}
		}
		return true
	})
	section := doc.Find(".content .pagination .next-page a")
	_, boolType := section.Attr("href")
	hasNext <- boolType
}

/**
 * 创建文件夹
 */
func createFolder(folderName string) (bool, error) {
	checkFolderNotExists, err := checkPathIsNotExists(folderName)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if checkFolderNotExists {
		err := os.MkdirAll(folderName, 0777)
		if err != nil {
			return false, err
		}
		fmt.Printf("create floder %s successful\n", folderName)
		return true, nil
	}
	return false, err
}

/**
 * 检查文件是否存在
 * 返回true 不存在， false 存在
 */
func checkPathIsNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//文章获取详情的分页
func fetchDetail(url string, savePath string, pagineCount chan int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	//分页计数
	detailCount := 0
	imgCount := make(chan int)
	doc.Find(".article-paging a").EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		detailURL, _ := contentSelection.Attr("href")
		filename := strconv.Itoa(i)
		go fetchImage(detailURL, savePath, filename, imgCount)
		imgFinish := <-imgCount
		fmt.Printf("finish this %d\n", imgFinish)
		detailCount++
		return true
	})
	pagineCount <- detailCount
}

//获取文章分页中的图片并进行下载操作
func fetchImage(url string, folderPath string, fileName string, count chan int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	content := doc.Find(".article-content img")
	//页面中的图片计数
	ct := 0
	content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		imgURL, _ := contentSelection.Attr("src")
		nFileName := strconv.Itoa(i) + "_" + fileName
		splitPoint := strings.Split(imgURL, ".")
		fileType := splitPoint[len(splitPoint)-1]
		go fetchAImage(imgURL, folderPath, nFileName, fileType, count)
		fmt.Printf("this page have %d images\n", <-count)
		ct++
		return true
	})
	count <- ct
}

//下载图片
func fetchAImage(img string, folderPath string, fileName string, fileType string, count chan int) {
	ct := 0
	respImg, err := http.Get(img)
	if err != nil {
		panic(err)
	}
	defer respImg.Body.Close()
	imgByte, _ := ioutil.ReadAll(respImg.Body)
	notExist, _ := checkPathIsNotExists(folderPath + "/" + fileName + "." + fileType)
	if notExist {
		fp, _ := os.Create(folderPath + "/" + fileName + "." + fileType)
		defer fp.Close()
		fp.Write(imgByte)
		ct++
		fmt.Printf("create %s/%s.%s successful\n", folderPath, fileName, fileType)
	}
	count <- ct
}
