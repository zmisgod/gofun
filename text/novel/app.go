package novel

import (
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

func getTextFile(textName string) string {
	return fmt.Sprintf("%s.txt", textName)
}

func FetchChapter(textName, url string) string {
	document, err := goquery.NewDocument(url)
	utils.CheckError(err)
	content := document.Find(".box_con3 #list dl dd")
	decoder := mahonia.NewDecoder("gbk")
	for j := 0; j <= 1125; j++ {
		time.Sleep(time.Duration(3) * time.Second)
		content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
			link, _ := contentSelection.Find("a").Attr("href")
			chapterName, _ := contentSelection.Find("a").Html()
			chapterName = decoder.ConvertString(chapterName)
			chatp := "第" + strconv.Itoa(j) + "章"
			fileName := textName+"/"+chatp+".txt"
			utils.CreateFile(fileName)
			file, err := os.OpenFile(getTextFile(textName), os.O_APPEND|os.O_WRONLY, 0777)
			utils.CheckError(err)
			if strings.Contains(chapterName, chatp) {
				go fetchFromServer(url+"/"+link, chapterName, file)
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
			log.Println("unexpected error")
			os.Exit(0)
		}
	} else {
		content = "\n\n" + chapterName + "\n\n" + content
		io.WriteString(fp, content)
		log.Printf("%s download finish\n", chapterName)
	}
}
