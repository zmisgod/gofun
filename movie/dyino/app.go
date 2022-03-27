package dyino

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"io/ioutil"
	"time"
)

type PublishVersionPara string

const (
	PublishVersionIMax       PublishVersionPara = "imax"
	PublishVersionIMaxStereo PublishVersionPara = "imaxStereo"
	Limit                    int                = 100
)

type ResponseData struct {
	Status   string       `json:"status"`
	Data     []*MovieInfo `json:"data"`
	Pageable struct {
		TotalElements uint `json:"totalElements"`
		TotalPages    uint `json:"totalPages"`
	} `json:"pageable"`
}

type MovieInfo struct {
	ID                 uint   `json:"id"`
	Code               string `json:"code"`
	Name               string `json:"name"`
	PublishVersion     string `json:"publishVersion"`
	PublishVersionName string `json:"publishVersionName"`
	Type               string `json:"type"`
	TypeName           string `json:"typeName"`
	ProducerName       string `json:"producerName"`
	PublisherName      string `json:"publisherName"`
	PublishDate        uint64 `json:"publishDate"`
}

func Fetch(ctx context.Context, year int, publishVersion PublishVersionPara, limit int) ([]*MovieInfo, error) {
	_url := fmt.Sprintf("https://api.zgdypw.cn/bits/w/porsfilms/api/s?page=0&size=%d&s_showYear=%d&s_publishVersion=%s&sort=id,desc", limit, year, publishVersion)
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodPost, _url, time.Duration(5)*time.Second, "", nil, utils.DefaultUserAgent)
	if err != nil {
		return []*MovieInfo{}, err
	}
	defer resp.Body.Close()
	return parseBody(ctx, resp.Body)
}

func parseBody(ctx context.Context, body io.Reader) ([]*MovieInfo, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return []*MovieInfo{}, err
	}
	var respData ResponseData
	if err := json.Unmarshal(data, &respData); err != nil {
		return []*MovieInfo{}, err
	}
	if respData.Status == "success" {
		return respData.Data, nil
	}
	return []*MovieInfo{}, errors.New("response error, status is " + respData.Status)
}
