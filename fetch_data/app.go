package fetch_data

import (
	"context"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"net/http"
)

type FetchData interface {
	Fetch(ctx context.Context) error
	ParseBody(ctx context.Context) error
	ToData(ctx context.Context) error
}

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

type Body struct {
	fetchBody  *http.Response
	HttpConfig HttpConfig
	DouBanId   string
}

type HttpConfig struct {
	Method  utils.HttpClientMethod
	Url     string
	Timeout int
	Proxy   string
	Body    io.ReadCloser
	Header  map[string]string
}

func NewFetch() (*Body, error) {
	douBanId := "1"
	a := Body{
		DouBanId: douBanId,
		HttpConfig: HttpConfig{
			Url: fmt.Sprintf("https://movie.douban.com/subject/%s/?from=subject-page", douBanId),
		},
	}
	return &a, nil
}

func Handle(ctx context.Context, fetch FetchData) error {
	if err := fetch.Fetch(ctx); err != nil {
		return err
	}
	if err := fetch.ParseBody(ctx); err != nil {
		return err
	}
	if err := fetch.ToData(ctx); err != nil {
		return err
	}
	return nil
}

func (a *Body) ParseBody(ctx context.Context) error {
	return nil
}

func (a *Body) ToData(ctx context.Context) error {
	return nil
}

func (a *Body) Fetch(ctx context.Context) error {
	var err error
	a.fetchBody, err = utils.HttpClient(ctx, a.HttpConfig.Method, a.HttpConfig.Url, a.HttpConfig.Timeout, a.HttpConfig.Proxy, a.HttpConfig.Body, a.HttpConfig.Header)
	if err != nil {
		return err
	}
	return nil
}

func parseBody(ctx context.Context, res interface{}) (string, error) {
	return "", nil
}
