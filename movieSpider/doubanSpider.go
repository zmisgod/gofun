package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var dbUser = "root"
var dbPass = "111111"
var dbTable = "mytest"
var targetURL = "https://movie.douban.com/j/search_subjects?type=movie&tag=%E7%83%AD%E9%97%A8&sort=recommend&page_limit=20&page_start="

type DoubanData struct {
	Subjects []detailData `json:"subjects"`
}
type detailData struct {
	Rate     string `json:"rate"`
	Cover_x  int    `json:"cover_x"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Playable bool   `json:"playable"`
	Cover    string `json:"cover"`
	Id       string `json:"id"`
	Cover_y  int    `json:"cover_y"`
	Is_new   bool   `json:"is_new"`
}

func main() {
	var isEnd bool
	startPage := 20
	for {
		resEnd := make(chan bool)
		go fetchURLData(targetURL+strconv.Itoa(startPage), resEnd)
		isEnd = <-resEnd
		if isEnd == false {
			break
		} else {
			startPage += 20
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
	var data DoubanData
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbTable+"?charset=utf8")
	checkErr(err)
	defer db.Close()
	for _, v := range data.Subjects {
		stmt, err := db.Prepare("INSERT ignore INTO douban_movie (id, rate, cover, title, url, playable,cover_x, cover_y, is_new) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		checkErr(err)

		_, err = stmt.Exec(v.Id, v.Rate, v.Cover, v.Title, v.Url, v.Playable, v.Cover_x, v.Cover_y, v.Is_new)
		checkErr(err)
	}
	end <- true
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
