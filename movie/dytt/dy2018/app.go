package dy2018

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/zmisgod/gofun/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func FetchByID(ctx context.Context, _url string) ([]string, error) {
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

func ReplaceUrl(_url string, replaceStr string) string {
	_urlInfo, err := url.Parse(_url)
	if err != nil {
		return ""
	}
	if replaceStr == "" {
		return _url
	}
	if replaceStr[0] != '/' {
		exp := strings.Split(_url, "/")
		if len(exp) > 0 {
			exp[len(exp)-1] = replaceStr
			return strings.Join(exp, "/")
		}
		return _url
	}
	return _urlInfo.Scheme + "://" + _urlInfo.Host + replaceStr
}

func SearchMovies(ctx context.Context, movieName string) (*SearchResult, error) {
	var result SearchResult
	result.SearchUrl = "https://www.dy2018.com/e/search/index.php"
	res := []byte(fmt.Sprintf("show=%s&tempid=1&keyboard=%s&Submit=%s&classid=0", utf2gbkEncode("title,smalltext"), utf2gbkEncode(movieName), utf2gbkEncode("立即搜索")))
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	trans, err := http.NewRequest("POST", result.SearchUrl, bytes.NewReader(res))
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
	result.LocationUrl = resp.Header.Get("Location")
	_targetUrl := ReplaceUrl(result.SearchUrl, result.LocationUrl)
	result.TargetUrl = _targetUrl
	result.SearchName = movieName
	return &result, err
}

type SearchResult struct {
	TotalPage   int    `json:"totalPage"`
	NowPage     int    `json:"nowPage"`
	LocationUrl string `json:"locationUrl"`
	TargetUrl   string `json:"targetUrl"`
	SearchName  string `json:"searchName"`
	SearchUrl   string `json:"searchUrl"`
}

func (a *SearchResult) HasMore(ctx context.Context) bool {
	if a.NowPage == 0 {
		return true
	}
	if a.NowPage > a.TotalPage {
		return false
	}
	return true
}

func (a *SearchResult) Next(ctx context.Context) ([]*SearchItem, error) {
	list, err := a.ListPage(ctx, a.NowPage)
	a.NowPage += 1
	return list, err
}

func (a *SearchResult) ListPage(ctx context.Context, page int) ([]*SearchItem, error) {
	a.NowPage = page
	_linkUrl := func(_url string, page int) string {
		if page == 0 {
			return _url
		}
		return fmt.Sprintf("%s&page=%d", _url, page)
	}
	_nowUrl := _linkUrl(a.TargetUrl, page)
	resp, err := http.Get(_nowUrl)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return a.parseListBody(string(bodyBytes))
}

type SearchItem struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func (a SearchItem) GetDownloadUrls(ctx context.Context) ([]string, error) {
	return FetchByID(ctx, a.Url)
}

func (a *SearchResult) parseListBody(bodyStr string) ([]*SearchItem, error) {
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
			Url:   ReplaceUrl(a.SearchUrl, downloadURL),
			Title: gb23122utf(title),
		})
		return true
	})
	totalPage := 0
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
	if a.TotalPage == 0 {
		a.TotalPage = totalPage
	}
	return list, nil
}

var pageReg = regexp.MustCompile(`(?m)page=(\d+)`)
