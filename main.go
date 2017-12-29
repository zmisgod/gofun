package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocument("http://www.alfuli.com/fuliba/page/14")
	if err != nil {
		log.Fatal(err)
	}
	count := make(chan int)
	doc.Find(".excerpt").EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		title := contentSelection.Find("h2 a").Text()
		targetURL, _ := contentSelection.Find("h2 a").Attr("href")
		folderName := strings.Replace(title, " ", "", -1)
		folder := "./images/" + folderName
		createFolder, err := createFolder(folder)
		if err != nil {
			panic(err)
		}
		if createFolder {
			go fetchDetail(targetURL, folder, count)
			allFinish := <-count
			fmt.Printf("title : %s download %d images \n", folderName, allFinish)
		} else {
			fmt.Printf("skip %s because exists\n", folder)
		}
		return true
	})
}

/**
 * 创建文件夹
 */
func createFolder(floderName string) (bool, error) {
	checkFloderNotExists, err := checkPathIsNotExists(floderName)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if checkFloderNotExists {
		err := os.MkdirAll(floderName, 0777)
		if err != nil {
			return false, err
		}
		fmt.Printf("create floder %s successful\n", floderName)
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
		imgurl, _ := contentSelection.Attr("src")
		nfilename := strconv.Itoa(i) + "_" + fileName
		splitPoint := strings.Split(imgurl, ".")
		fileType := splitPoint[len(splitPoint)-1]
		go fetchAImage(imgurl, folderPath, nfilename, fileType, count)
		thisPageImage := <-count
		fmt.Printf("this page have %d images\n", thisPageImage)
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
