package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

var start = "100551"

func main() {
	for i := 0; i < 10; i++ {
		go test(i)
	}
}

func test(id int) {
	createFloder(strconv.Itoa(id))
}

func createFloder(name string) {
	nowTime := int(time.Now().Unix())
	err := os.Mkdir(name+"_"+strconv.Itoa(nowTime), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func fetch(start string) {
	rest, err := getMovieURL("https://www.dy2018.com/i/" + start + ".html")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(rest)
	}
}

func getMovieURL(baiduLink string) ([]string, error) {
	doc, err := goquery.NewDocument(baiduLink)
	if err != nil {
		return []string{}, err
	}
	content := doc.Find("#Zoom table")

	downloads := []string{}
	content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		downloadURL, exists := contentSelection.Find("tbody tr td a").Attr("href")
		if !exists {
			return false
		}
		enc := mahonia.NewDecoder("gbk")
		utf := enc.ConvertString(downloadURL)
		downloads = append(downloads, utf)
		return true
	})
	return downloads, nil
}
