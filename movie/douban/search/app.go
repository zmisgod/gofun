package search

import (
	"context"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type DouBanMovieInfo struct {
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Cover       string `json:"cover"`
	Rate        string `json:"rate"`
	Introduce   string `json:"introduce"`
	IMDB        string `json:"imdb"`
	Minute      int    `json:"minute"`
	ReleaseDate string `json:"release_date"`
}

func Fetch(ctx context.Context, douBanId string) (*DouBanMovieInfo, error) {
	_url := fmt.Sprintf("https://movie.douban.com/subject/%s/?from=subject-page", douBanId)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, 5, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func parseBody(ctx context.Context, respBody io.Reader) (*DouBanMovieInfo, error) {
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
			preName, _ := GetValueByPrefix(ctx, strings.NewReader(v[1]), "span")
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
	dataObj := &DouBanMovieInfo{
		Name:        name,
		Year:        year,
		Cover:       cover,
		Rate:        rate,
		Introduce:   strings.TrimSpace(introduce),
		IMDB:        imdbString,
		Minute:      min,
		ReleaseDate: releaseDate,
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
			_n, _ := GetValueByPrefix(ctx, strings.NewReader(j[1]), "a")
			if _n != "" {
				return _n
			}
		}
	}
	return ""
}

func GetValueByPrefix(ctx context.Context, body io.Reader, prefixS string) (string, error) {
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
