package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"time"
)

var targetURL = "http://www.xiuren.org/tuigirl-001.html"

func main() {
	proxyStr, err := url.Parse("http://127.0.0.1:1087")
	if err != nil {
		fmt.Println(err)
	}else {
		tr := &http.Transport{
			Proxy: http.ProxyURL(proxyStr),
		}
		client := &http.Client{
			Transport: tr,
			Timeout:   time.Second * 10, //超时时间
		}
		request, err := http.NewRequest("GET", targetURL, nil)
		if err == nil {
			request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36")
			resp, err := client.Do(request)
			if err == nil {
				defer resp.Body.Close()
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err == nil {
					sec := doc.Find("#main #post .post .photoThum")
					sec.EachWithBreak(func(i int, selection *goquery.Selection) bool {
						imagesURL, exist := selection.Find(" a").Attr("href")
						if exist {
							fmt.Println(imagesURL)
						}
						return true
					})
				}
			}
		}
	}
}
