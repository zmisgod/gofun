package imdb

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/zmisgod/gofun/utils"
	"io"
	"strings"
	"time"
)

type MovieInfo struct {
	Name                   string `json:"name"`
	RunTime                string `json:"run_time"`
	SoundMix               string `json:"sound_mix"`
	Color                  string `json:"color"`
	AspectRatio            string `json:"aspect_ratio"`
	Camera                 string `json:"camera"`
	Laboratory             string `json:"laboratory"`
	FilmLength             string `json:"film_length"`
	NegativeFormat         string `json:"negative_format"`
	CinematographicProcess string `json:"cinematographic_process"`
	PrintedFilmFormat      string `json:"printed_film_format"`
}

func Fetch(ctx context.Context, imdbStr string) (*MovieInfo, error) {
	_url := fmt.Sprintf("https://www.imdb.com/title/%s/technical?ref_=tt_dt_spec", imdbStr)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func parseBody(ctx context.Context, body io.Reader) (*MovieInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}
	mainSec := doc.Find("#main")
	titleStr, err := mainSec.Find(".subpage_title_block__right-column h3 a").Html()
	var movieInfo MovieInfo
	movieInfo.Name = titleStr

	techSec := mainSec.Find("#technical_content tr")

	resMap := make(map[string]string, 0)
	configList := make([][]string, 0)
	techSec.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		tdSelection := selection.Find("td")
		oneConfigList := make([]string, 2)
		tdSelection.EachWithBreak(func(j int, se *goquery.Selection) bool {
			name, _ := se.Html()
			if j == 0 {
				oneConfigList[0] = strings.TrimSpace(name)
			} else if j == 1 {
				oneConfigList[1] = strings.TrimSpace(name)
			}
			return true
		})
		configList = append(configList, oneConfigList)
		return true
	})
	for _, v := range configList {
		if len(v) == 2 {
			resMap[v[0]] = v[1]
		}
	}
	for k, v := range resMap {
		if k == "Runtime" {
			movieInfo.RunTime = v
		} else if k == "Sound Mix" {
			movieInfo.SoundMix = v
		} else if k == "Color" {
			movieInfo.Color = v
		} else if k == "Aspect Ratio" {
			movieInfo.AspectRatio = v
		} else if k == "Camera" {
			movieInfo.Camera = v
		} else if k == "Laboratory" {
			movieInfo.Laboratory = v
		} else if k == "Film Length" {
			movieInfo.FilmLength = v
		} else if k == "Negative Format" {
			movieInfo.NegativeFormat = v
		} else if k == "Cinematographic Process" {
			movieInfo.CinematographicProcess = v
		} else if k == "Printed Film Format" {
			movieInfo.PrintedFilmFormat = v
		}
	}
	return &movieInfo, nil
}
