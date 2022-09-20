package dy2018

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/zmisgod/gofun/utils"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FetchByID(ctx context.Context, id string) ([]string, error) {
	_url := fmt.Sprintf("https://www.dy2018.com/i/%s.html", id)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	return parseDetailBody(ctx, resp.Body)
}

func FetchByUrl(ctx context.Context, _url string) ([]string, error) {
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()
	return parseDetailBody(ctx, resp.Body)
}

func parseDetailBody(ctx context.Context, body io.Reader) ([]string, error) {
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

func gb23122utf(str string) string {
	downloadURL, _ := url.PathUnescape(str)
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(downloadURL)
}

func utf2gbkEncode(str string) string {
	enc := mahonia.NewEncoder("gbk")
	return enc.ConvertString(str)
}

func SearchMovies(ctx context.Context, movieName string) (*SearchResult, error) {
	res := []byte(fmt.Sprintf("show=%s&tempid=1&keyboard=%s&Submit=%s&classid=0", utf2gbkEncode("title,smalltext"), utf2gbkEncode(movieName), utf2gbkEncode("立即搜索")))
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	trans, err := http.NewRequest("POST", "https://www.dy2018.com/e/search/index.php", bytes.NewReader(res))
	if err != nil {
		return nil, err
	}
	trans.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(trans)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var result SearchResult
	result.LocationUrl = resp.Header.Get("Location")
	result.SearchName = movieName
	return nil, err
}

type SearchResult struct {
	List        []*SearchItem `json:"list"`
	Page        int           `json:"page"`
	LocationUrl string        `json:"locationUrl"`
	SearchName  string        `json:"searchName"`
}

func (a SearchResult) ListPage(ctx context.Context, page uint) ([]*SearchItem, error) {
	list := make([]*SearchItem, 0)
	return list, nil
}

type SearchItem struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func parseListBody(bodyStr string) (*SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}
	content := doc.Find(".co_content8 ul table")
	list := make([]*SearchItem, 0)
	content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		downloadURL, linkExists := contentSelection.Find("tbody tr td a").Attr("href")
		title, titleExists := contentSelection.Find("tbody tr td a").Attr("title")
		if !linkExists && !titleExists {
			return false
		}
		list = append(list, &SearchItem{
			Url:   downloadURL,
			Title: gb23122utf(title),
		})
		return true
	})
	totalPage := 1
	pagination := doc.Find(".co_content8 .x a")
	pagination.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		res, _ := contentSelection.Html()
		if gb23122utf(res) == "尾页" {
			_url, _ := contentSelection.Attr("href")
			exp := pageReg.FindStringSubmatch(_url)
			if len(exp) >= 2 {
				totalPage, _ = strconv.Atoi(exp[1])
			}
		}
		return true
	})
	return &SearchResult{
		List: list,
		Page: totalPage,
	}, nil
}

var pageReg = regexp.MustCompile(`(?m)page=(\d+)`)
