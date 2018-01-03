package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
type PostJSON struct {
	Code int          `json:"code"`
	Data []detailJSON `json:"data"`
	Msg  string       `json:"msg"`
}
type detailJSON struct {
	Id         int    `json:"id"`
	Post_date  string `json:"post_data"`
	Post_intro string `json:"post_intro"`
	Post_title string `json:"post_title"`
}

func main() {
	url := "https://movie.douban.com/j/search_subjects?type=movie&tag=%E8%B1%86%E7%93%A3%E9%AB%98%E5%88%86&sort=rank&page_limit=20&page_start=40"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}
	var data DoubanData
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
