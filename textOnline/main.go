package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

//URL target
var URL = "https://www.18xs.org/book_25306/"

var textName = "丹武双绝.txt"

func main() {
	createFile()
	fetchChapter(URL)
}

func fetchChapter(url string) string {
	document, err := goquery.NewDocument(url)
	checkError(err)
	content := document.Find(".box_con3 #list dl dd")
	decoder := mahonia.NewDecoder("gbk")
	file, err := os.OpenFile(textName, os.O_APPEND|os.O_WRONLY, 0777)
	checkError(err)
	for j := 0; j <= 1125; j++ {
		time.Sleep(time.Duration(3) * time.Second)
		content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
			link, _ := contentSelection.Find("a").Attr("href")
			chapterName, _ := contentSelection.Find("a").Html()
			chapterName = decoder.ConvertString(chapterName)
			chatp := "第" + strconv.Itoa(j) + "章"
			fmt.Println(chatp)
			if strings.Contains(chapterName, chatp) {
				fetchFromServer(URL+"/"+link, chapterName, file)
				fmt.Printf("%s download finish\n", chapterName)
			}
			return true
		})
	}
	return ""
}

func fetchFromServer(url, chapterName string, fp *os.File) {
	document, err := goquery.NewDocument(url)
	checkError(err)
	content, err := document.Find("#content").Html()
	checkError(err)
	decoder := mahonia.NewDecoder("gbk")
	content = strings.Replace(content, "&nbsp;", "", len(content))
	content = strings.Replace(content, "<br/>", "", len(content))
	content = decoder.ConvertString(content)
	content = strings.Replace(content, "聽", "", len(content))
	if content == "" {
		unavailableStr, _ := document.Find("title").Html()
		//如果遇到服务不可用503，直接等待10秒后重试
		if strings.Contains(unavailableStr, "503") {
			time.Sleep(time.Duration(1) * time.Second)
			fetchFromServer(url, chapterName, fp)
		} else {
			fmt.Println("unexpected error")
			os.Exit(0)
		}
	} else {
		content = "\n\n" + chapterName + "\n\n" + content
		io.WriteString(fp, content)
	}
}

func createFile() {
	_, err := os.Stat(textName)
	if err != nil {
		file, err := os.Create(textName)
		checkError(err)
		defer file.Close()
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
