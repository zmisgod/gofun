package search

import (
	"context"
	"errors"
	"github.com/zmisgod/gofun/movie/dytt/dy2018"
	"github.com/zmisgod/gofun/utils"
	_ "io/ioutil"
	_ "net/http"
	"net/url"
	_ "net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func Fetch(ctx context.Context, movieName string) ([]string, error) {
	baiduURL := "http://www.baidu.com/s?wd=site%3Ady2018.com%20" + url.PathEscape(movieName)
	baiduLink, err := searchFromBaidu(ctx, baiduURL)
	if err != nil {
		return []string{}, err
	}
	res, err := dy2018.FetchByUrl(ctx, baiduLink)
	if err == nil {
		return res, nil
	}
	return []string{}, err
}

func getDyId(ctx context.Context, targetUrl string) string {
	var re = regexp.MustCompile(`(?m)\/([0-9]*).html`)
	list := re.FindAllStringSubmatch(targetUrl, 10000)
	if len(list) > 0 && len(list[0]) > 1 {
		return list[0][1]
	}
	return ""
}

func searchFromBaidu(ctx context.Context, baiduURL string) (string, error) {
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, baiduURL, 5, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
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
