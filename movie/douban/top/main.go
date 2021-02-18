package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbTable = "douban_year_best_movie"
var pageLimit = "20"
var doubanCategoory = []string{"热门", "最新", "经典", "可播放", "豆瓣高分", "冷门佳片", "华语", "欧美", "日本", "动作", "喜剧", "爱情", "科幻", "悬疑", "恐怖", "成长"}

//DoubanData data
type DoubanData struct {
	Subjects []detailData `json:"subjects"`
}

type detailData struct {
	Rate     string `json:"rate"`
	CoverX   int    `json:"cover_x"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Playable bool   `json:"playable"`
	Cover    string `json:"cover"`
	ID       string `json:"id"`
	CoverY   int    `json:"cover_y"`
	IsNew    bool   `json:"is_new"`
}

var dbCon *sql.DB

func main() {
	err := godotenv.Load("./../../.env")
	checkError(err)
	dbHost := os.Getenv("mysql.host")
	dbPort := os.Getenv("mysql.port")
	dbUser := os.Getenv("mysql.username")
	dbPass := os.Getenv("mysql.password")
	dbName := os.Getenv("mysql.dbname")
	dbCon, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	checkError(err)
	defer dbCon.Close()
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func goRunOneTime(targetURL string, next chan bool) {
	var isEnd int
	startPage := 0
	for {
		resEnd := make(chan int)
		//一秒钟后请求，防止豆瓣接口屏蔽ip
		time.Sleep(time.Second * 3)
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
	i := 1
	for _, v := range data.Subjects {

		stmt, err := dbCon.Prepare("INSERT ignore INTO " + dbTable + " (id, rate, cover, title, url, playable,cover_x, cover_y, is_new) values (?, ?, ?, ?, ?, ?, ?, ?, ?)")
		defer stmt.Close()
		checkErr(err)

		_, err = stmt.Exec(v.ID, v.Rate, v.Cover, v.Title, v.URL, v.Playable, v.CoverX, v.CoverY, v.IsNew)
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
