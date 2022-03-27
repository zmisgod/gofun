package name

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io/ioutil"
	"time"
)

type SearchDouBanInfo struct {
	Episode  string `json:"episode"`
	Id       string `json:"id"`
	Img      string `json:"img"`
	SubTitle string `json:"sub_title"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Url      string `json:"url"`
	Year     string `json:"year"`
}

func Fetch(ctx context.Context, movieName string) ([]*SearchDouBanInfo, error) {
	_url := fmt.Sprintf("https://movie.douban.com/j/subject_suggest?q=%s", movieName)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var one []*SearchDouBanInfo
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return one, err
	}
	err = json.Unmarshal(resBody, &one)
	return one, err
}
