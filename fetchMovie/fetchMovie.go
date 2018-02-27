package main 

import (
	"os"
	"fmt"
	"flag"
	"net/url"
	"github.com/axgle/mahonia"
	_ "io/ioutil"
	_ "net/http"
	"errors"
	"github.com/PuerkitoBio/goquery"
	_ "net/url"
)

func main () {
	if len(os.Args) <= 1 {
		fmt.Println("please input movive")
		os.Exit(0)
	}
	movieName := flag.String("movieName", os.Args[1], "sasasa")
	flag.Parse()
	baiduURL := "http://www.baidu.com/s?wd=" + url.PathEscape(*movieName) + "+%E7%94%B5%E5%BD%B1%E5%A4%A9%E5%A0%82"
	baiduLink, err := searchFromBaidu(baiduURL)
	 if err != nil {
		panic(err)
	}
	res, err := getMovieURL(baiduLink)
	if err != nil {
		panic(err)
	}
	for _,v := range(res) {
		fmt.Println(v)
	}
}

func searchFromBaidu(baiduURL string) (string, error){
	doc, err := goquery.NewDocument(baiduURL)
	if err != nil {
		return "", err
	}
	section := doc.Find("#content_left #1 .t a")
	nextURL, exists := section.Attr("href")
	if exists {
		return nextURL, nil
	}
	return "", errors.New("没找到")
}

func getMovieURL(baiduLink string) ([]string, error){
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