package welfare

import (
	"encoding/json"
	"flag"
	"github.com/zmisgod/gofun/downloader"
	"github.com/zmisgod/gofun/utils"
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
		NewSpider(*startPage, *count)
	}
}

func NewSpider(startPage, count int) {
	for i := startPage; i < count; i++ {
		haveNext := make(chan bool)
		nowURL := targetURL + strconv.Itoa(i)
		log.Println(nowURL)
		go fetchPage(nowURL, haveNext)
		if !<-haveNext {
			log.Println("this page has no more content to fetch")
		}
	}
}

type GenJson struct {
	Title string   `json:"title"`
	Files []string `json:"files"`
}

func GenerateJson() {
	file, err := ioutil.ReadDir("images")
	utils.CheckError(err)
	var jsons []GenJson
	for _, v := range file {
		var oneJson GenJson
		if v.IsDir() && v.Name() != "" {
			oneJson.Title = v.Name()
			image, err := ioutil.ReadDir("images/" + v.Name())
			if err != nil {
				log.Println(err)
			} else {
				for _, j := range image {
					if !j.IsDir() && j.Name() != "" {
						oneJson.Files = append(oneJson.Files, j.Name())
					}
				}
				if len(oneJson.Files) == 0 {
					log.Println("images/" + v.Name())
				}
			}
		}
		jsons = append(jsons, oneJson)
	}
	writeData, err := json.Marshal(jsons)
	utils.CheckError(err)
	err = ioutil.WriteFile("data.json", writeData, 0777)
	utils.CheckError(err)
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

//func fetchNextURL(nowURL string, chanNextURL chan string) {
//	doc, err := goquery.NewDocument(nowURL)
//	if err != nil {
//		panic(err)
//	}
//	section := doc.Find(".content .pagination .next-page a")
//	nextURL, _ := section.Attr("href")
//	log.Println("next url = " + nextURL + "\n")
//	chanNextURL <- nextURL
//}

func fetchPage(url string, hasNext chan bool) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
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
			createFolder, err := utils.CreateFolder(folder)
			if err != nil {
				return true
			}
			if createFolder {
				go fetchDetail(targetURL, folder, count)
				ct++
				allFinish := <-count
				log.Printf("title : %s download %d images \n", folderName, allFinish)
			} else {
				log.Printf("skip %s because exists\n", folder)
			}
		}
		return true
	})
	section := doc.Find(".content .pagination .next-page a")
	_, boolType := section.Attr("href")
	hasNext <- boolType
}

//文章获取详情的分页
func fetchDetail(url string, savePath string, pagineCount chan int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	//分页计数
	detailCount := 0
	imgCount := make(chan int)
	doc.Find(".article-paging a").EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		detailURL, _ := contentSelection.Attr("href")
		filename := strconv.Itoa(i)
		go fetchImage(detailURL, savePath, filename, imgCount)
		imgFinish := <-imgCount
		log.Printf("finish this %d\n", imgFinish)
		detailCount++
		return true
	})
	pagineCount <- detailCount
}

//获取文章分页中的图片并进行下载操作
func fetchImage(url string, folderPath string, fileName string, count chan int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
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
		log.Printf("this page have %d images\n", <-count)
		ct++
		return true
	})
	count <- ct
}

//下载图片
func fetchAImage(img string, folderPath string, fileName string, fileType string, count chan int) {
	ct := 0
	res, err := downloader.NewDownloader(img)
	if err != nil {
		log.Println(err)
	} else {
		res.SetSavePath(folderPath)
		res.SetSaveName(fileName + "." + fileType)
		err = res.SaveFile()
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("create %s/%s.%s successful\n", folderPath, fileName, fileType)
			ct++
		}
	}
	count <- ct
}
