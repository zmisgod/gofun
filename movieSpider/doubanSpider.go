package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dbUser = "root"
var dbPass = "111111"
var dbName = "mytest"
var dbTable = "douban_year_best_movie"
var pageLimit = "20"
var doubanCategoory = []string{"热门", "最新", "经典", "可播放", "豆瓣高分", "冷门佳片", "华语", "欧美", "日本", "动作", "喜剧", "爱情", "科幻", "悬疑", "恐怖", "成长"}

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
	for _, v := range doubanCategoory {
		var tem = "https://movie.douban.com/j/search_subjects?type=movie&tag=" + url.QueryEscape(v) + "&sort=recommend&page_limit=" + pageLimit + "&page_start="
		continueNext := make(chan bool)
		go goRunOneTime(tem, continueNext)
		cN := <-continueNext
		if cN {
			fmt.Println("continue next fetch")
		}
	}
}

func goRunOneTime(targetURL string, next chan bool) {
	var isEnd int
	startPage := 0
	for {
		resEnd := make(chan int)
		//一秒钟后请求，防止豆瓣接口屏蔽ip
		time.Sleep(time.Second * 1)
		go fetchURLData(targetURL+strconv.Itoa(startPage), resEnd)
		isEnd = <-resEnd
		if isEnd == 0 {
			break
		} else {
			startPage += 20
		}
	}
	next <- true
}

func fetchURLData(url string, end chan int) {
	resp, err := http.Get(url)
	if err != nil {
		end <- 0
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		end <- 0
	}
	var data DoubanData
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		end <- 0
	}
	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8")
	if err != nil {
		end <- 0
	}
	defer db.Close()
	i := 1
	for _, v := range data.Subjects {

		stmt, err := db.Prepare("INSERT ignore INTO " + dbTable + " (id, rate, cover, title, url, playable,cover_x, cover_y, is_new) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		checkErr(err)

		_, err = stmt.Exec(v.Id, v.Rate, v.Cover, v.Title, v.Url, v.Playable, v.Cover_x, v.Cover_y, v.Is_new)
		if err == nil {
			i++
		}
	}
	end <- i
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
