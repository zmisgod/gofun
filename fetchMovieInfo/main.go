package main

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var urls = "http://www.baidu.com/s?wd=三块广告牌+豆瓣电影"

func main() {
	fetchURL, err := searchFromBaidu(urls)
	if err != nil {
		panic(err)
	}
	findMovieInfo(fetchURL)
}

func searchFromBaidu(baiduURL string) (string, error) {
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

type movieInfo struct {
	Name          string `json:"name"`
	DoubanID      int64  `json:"douban_id"`
	MYear         string `json:"m_year"`
	Director      string `json:"director"`
	Screenwriter  string `json:"screenwriter"`
	Rate          string `json:"rate"`
	MType         string `json:"m_type"`
	Website       string `json:"website"`
	Country       string `json:"country"`
	Language      string `json:"language"`
	ReleaseDate   string `json:"release_date"`
	MLength       string `json:"m_length"`
	AlternateName string `json:"alternate_name"`
	Imdb          string `json:"imdb"`
	Synopsis      string `json:"synopsis "`
	MCover        string `json:"m_cover"`
}

var movie movieInfo

func findMovieInfo(doubanURL string) {
	doc, err := goquery.NewDocument(doubanURL)
	if err != nil {
		panic(err)
	}
	nameSection := doc.Find("#wrapper #content h1 span")
	name, err := nameSection.First().Html()
	if err != nil {
		panic(err)
	}
	movie.Name = name
	year, _ := nameSection.Next().Html()
	movie.MYear = year
	contentSection := doc.Find("#wrapper #content .clearfix .article .indent .subjectwrap .subject")
	cover, _ := contentSection.Find("#mainpic a img").Attr("src")
	rate := doc.Find("#wrapper #content .clearfix .article .indent .subjectwrap  #interest_sectl .rating_wrap .rating_self .rating_num").Text()

	fmt.Println(name)
	fmt.Println(year)
	fmt.Println(cover)
	fmt.Println(rate)
}
