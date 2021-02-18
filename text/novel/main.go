package main

import (
	"fmt"
	"github.com/zmisgod/goSpider/utils"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

func main() {
	var textName = "丹武双绝.txt"
	var url = "https://www.18xs.org/book_25306/"
	utils.CreateFile(textName)
	fetchChapter(textName, url)
}

func fetchChapter(textName, url string) string {
	document, err := goquery.NewDocument(url)
	utils.CheckError(err)
	content := document.Find(".box_con3 #list dl dd")
	decoder := mahonia.NewDecoder("gbk")
	file, err := os.OpenFile(textName, os.O_APPEND|os.O_WRONLY, 0777)
	utils.CheckError(err)
	for j := 0; j <= 1125; j++ {
		time.Sleep(time.Duration(3) * time.Second)
		content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
			link, _ := contentSelection.Find("a").Attr("href")
			chapterName, _ := contentSelection.Find("a").Html()
			chapterName = decoder.ConvertString(chapterName)
			chatp := "第" + strconv.Itoa(j) + "章"
			fmt.Println(chatp)
			if strings.Contains(chapterName, chatp) {
				fetchFromServer(url+"/"+link, chapterName, file)
				fmt.Printf("%s download finish\n", chapterName)
			}
			return true
		})
	}
	return ""
}

func fetchFromServer(url, chapterName string, fp *os.File) {
	document, err := goquery.NewDocument(url)
	utils.CheckError(err)
	content, err := document.Find("#content").Html()
	utils.CheckError(err)
	decoder := mahonia.NewDecoder("gbk")
	content = strings.Replace(content, "&nbsp;", "", len(content))
	content = strings.Replace(content, "<br/>", "", len(content))
	content = decoder.ConvertString(content)
	content = strings.Replace(content, "聽", "", len(content))
	if content == "" {
		unavailableStr, _ := document.Find("title").Html()
		//如果遇到服务不可用503，直接等待10秒后重试
		if strings.Contains(unavailableStr, "503") {
			time.Sleep(time.Second)
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
