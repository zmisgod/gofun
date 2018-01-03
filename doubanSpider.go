package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var targetURL = "https://movie.douban.com/j/search_subjects?type=movie&tag=%E7%83%AD%E9%97%A8&sort=recommend&page_limit=20&page_start="

func main() {
	end := false
	for {
		if end {
			break
		}
	}
}

func fetchURLData(url string, end chan bool) {
	resp, err := http.Get(url)
	if err != nil {
		end <- true
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		end <- true
	}
	fmt.Println(respBody)
	end <- false
}
