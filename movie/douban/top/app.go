package top

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"io/ioutil"
	"time"
)

type DouBanCategory string

const (
	DouBanCategoryHot      DouBanCategory = "热门"
	DouBanCategoryNew      DouBanCategory = "最新"
	DouBanCategoryClass    DouBanCategory = "经典"
	DouBanCategoryCanPlay  DouBanCategory = "可播放"
	DouBanCategoryTop      DouBanCategory = "豆瓣高分"
	DouBanCategoryCold     DouBanCategory = "冷门佳片"
	DouBanCategoryAsia     DouBanCategory = "华语"
	DouBanCategoryEurope   DouBanCategory = "欧美"
	DouBanCategoryJapan    DouBanCategory = "日本"
	DouBanCategoryAction   DouBanCategory = "动作"
	DouBanCategoryComedy   DouBanCategory = "喜剧"
	DouBanCategoryLove     DouBanCategory = "爱情"
	DouBanCategoryScience  DouBanCategory = "科幻"
	DouBanCategorySuspense DouBanCategory = "悬疑"
	DouBanCategoryTerror   DouBanCategory = "恐怖"
	DouBanCategoryGrow     DouBanCategory = "成长"
)

type DouBanData struct {
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

func Fetch(ctx context.Context, category DouBanCategory, startId int, pageSize int) (*DouBanData, error) {
	targetUrl :=
		fmt.Sprintf("https://movie.douban.com/j/search_subjects?type=movie&tag=%s&sort=recommend&page_limit=%d&page_start=%d",
			string(category), pageSize, startId)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, targetUrl, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func parseBody(ctx context.Context, body io.Reader) (*DouBanData, error) {
	respBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var data *DouBanData
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
