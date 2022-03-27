package id

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zmisgod/gofun/utils"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DouBanMovieInfo struct {
	Name        string   `json:"name"`
	Year        int      `json:"year"`
	Cover       string   `json:"cover"`
	Rate        string   `json:"rate"`
	Introduce   string   `json:"introduce"`
	IMDB        string   `json:"imdb"`
	Minute      int      `json:"minute"`
	ReleaseDate string   `json:"release_date"`
	CoverUrls   []string `json:"cover_urls"`
	DoubanId    string   `json:"douban_id"`
}

func Fetch(ctx context.Context, douBanId string) (*DouBanMovieInfo, error) {
	_url := fmt.Sprintf("https://movie.douban.com/subject/%s/?from=subject-page", douBanId)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, douBanId, resp.Body)
}

func parseBody(ctx context.Context, douBanId string, respBody io.Reader) (*DouBanMovieInfo, error) {
	doc, err := goquery.NewDocumentFromReader(respBody)
	if err != nil {
		return nil, err
	}
	nameSection := doc.Find("#wrapper #content h1 span")
	name, err := nameSection.First().Html()
	if err != nil {
		return nil, err
	}
	yearStr, _ := nameSection.Next().Html()
	contentSection := doc.Find("#wrapper #content .clearfix .article .indent .subjectwrap .subject")
	cover, _ := contentSection.Find("#mainpic a img").Attr("src")
	rate := doc.Find("#wrapper #content .clearfix .article .indent .subjectwrap  #interest_sectl .rating_wrap .rating_self .rating_num").Text()
	introduce, _ := doc.Find("#link-report span").Html()
	movieInfo, _ := doc.Find("#info").Html()
	var reB = regexp.MustCompile(`(?m)(.*)[<br>|<br\/>]`)
	infoArr := reB.FindAllStringSubmatch(movieInfo, 1000000)
	var imdbString string
	var _min string
	var releaseDate string
	for _, v := range infoArr {
		if len(v) > 1 {
			preName, _ := getValueByPrefix(ctx, strings.NewReader(v[1]), "span")
			if preName == "IMDb链接:" {
				imdbString = getImdb(ctx, v[1])
			} else if preName == "片长:" {
				_min = getChinaMinute(ctx, v[1])
				if _min == "" {
					_min = getMinute(ctx, v[1])
				}
			} else if preName == "上映日期:" {
				releaseDate = getReleaseDate(ctx, v[1])
			}
		}
	}

	var re = regexp.MustCompile(`(?m)[0-9]+`)
	var year int
	yearList := re.FindAllStringSubmatch(yearStr, 1000)
	if len(yearList) > 0 && len(yearList[0]) > 0 {
		year, _ = strconv.Atoi(yearList[0][0])
	}
	min, _ := strconv.Atoi(_min)
	time.Sleep(time.Duration(utils.Rand(3, 10)) * time.Second)
	coverUrls, _ := FetchMovieCover(ctx, douBanId)
	dataObj := &DouBanMovieInfo{
		Name:        name,
		Year:        year,
		Cover:       cover,
		Rate:        rate,
		Introduce:   strings.TrimSpace(introduce),
		IMDB:        imdbString,
		Minute:      min,
		ReleaseDate: releaseDate,
		CoverUrls:   coverUrls,
		DoubanId:    douBanId,
	}
	return dataObj, nil
}

func getChinaMinute(ctx context.Context, sBody string) string {
	var re = regexp.MustCompile(`(?m)([0-9]+)分钟\(中国大陆\)`)
	minAList := re.FindAllStringSubmatch(sBody, 100000)
	for _, j := range minAList {
		if len(j) > 1 {
			return j[1]
		}
	}
	return ""
}

func getMinute(ctx context.Context, sBody string) string {
	var re = regexp.MustCompile(`(?m)([0-9]+)分钟`)
	minAList := re.FindAllStringSubmatch(sBody, 100000)
	for _, j := range minAList {
		if len(j) > 1 {
			return j[1]
		}
	}
	return ""
}

func getReleaseDate(ctx context.Context, sBody string) string {
	var re = regexp.MustCompile(`(?m)([0-9]+\-[0-9]+\-[0-9]+)\(中国大陆\)`)
	minAList := re.FindAllStringSubmatch(sBody, 100000)
	for _, j := range minAList {
		if len(j) > 1 {
			return j[1]
		}
	}
	return ""
}

func getImdb(ctx context.Context, sBody string) string {
	reC := regexp.MustCompile(`(?m)<\/span>(.*)<`)
	infoAList := reC.FindAllStringSubmatch(sBody, 100000)
	for _, j := range infoAList {
		if len(j) > 1 {
			_n, _ := getValueByPrefix(ctx, strings.NewReader(j[1]), "a")
			if _n != "" {
				return _n
			}
		}
	}
	return ""
}

func getValueByPrefix(ctx context.Context, body io.Reader, prefixS string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}
	var section *goquery.Selection
	if prefixS != "" {
		section = doc.Find(prefixS)
	} else {
		section = doc.Selection
	}
	return section.Html()
}

func FetchMovieCover(ctx context.Context, douBanId string) ([]string, error) {
	_url := fmt.Sprintf("https://movie.douban.com/subject/%s/photos?type=R", douBanId)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseCoverBody(ctx, resp.Body)
}

func parseCoverBody(ctx context.Context, respBody io.Reader) ([]string, error) {
	result := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(respBody)
	if err != nil {
		return nil, err
	}
	liSection := doc.Find(".poster-col3 li")
	liSection.EachWithBreak(func(i int, liOneSelection *goquery.Selection) bool {
		oneImg, ex := liOneSelection.Find("img").Attr("src")
		if ex {
			result = append(result, strings.Replace(oneImg, "/m/", "/l/", 300))
		}
		return true
	})
	return result, nil
}
