package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func fetchNextURL(nowURL string, chanNextURL chan string) {
	doc, err := goquery.NewDocument(nowURL)
	if err != nil {
		panic(err)
	}
	section := doc.Find(".content .pagination .next-page a")
	nextURL, _ := section.Attr("href")
	fmt.Println(nextURL + "\n")
	chanNextURL <- nextURL
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(`please input at least one paramter
./main [target url (string)] [pull page(int), default:1, max:each all page]

such as : 
./main http://www.alfuli.com/fuliba/page/13
or
./main http://www.alfuli.com/fuliba/page/13 10000000
`)
		os.Exit(0)
	}
	targetURL := ""
	eachTime := 1
	for k, v := range os.Args {
		if k == 0 {
			continue
		}
		if k == 1 {
			targetURL = v
		}
		if k == 2 {
			eachTime, _ = strconv.Atoi(v)
		}
	}
	if targetURL == "" {
		fmt.Println("please input a target url instead of empty")
		os.Exit(0)
	}
	nowURL := targetURL
	for i := 0; i < eachTime; i++ {
		nextURL := make(chan string)
		go fetchNextURL(nowURL, nextURL)
		nowURL = <-nextURL
	}
}
