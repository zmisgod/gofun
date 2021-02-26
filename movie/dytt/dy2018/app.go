package dy2018

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/zmisgod/gofun/utils"
	"io"
)

func FetchByID(ctx context.Context, id string) ([]string, error) {
	_url := fmt.Sprintf("https://www.dy2018.com/i/%s.html", id)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, 5, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func FetchByUrl(ctx context.Context, _url string)([]string, error) {
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, 5, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func parseBody(ctx context.Context, body io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return []string{}, err
	}
	content := doc.Find("#Zoom table")
	downloads := make([]string, 0)
	content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		downloadURL, exists := contentSelection.Find("tbody tr td a").Attr("href")
		if !exists {
			return false
		}
		enc := mahonia.NewDecoder("gbk")
		utf := enc.ConvertString(downloadURL)
		downloads = append(downloads, utf)
		return true
	})
	return downloads, nil
}
