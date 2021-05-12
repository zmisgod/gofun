package plus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zmisgod/gofun/easy_http_client"
	"io/ioutil"
)

type RespData struct {
	Code uint           `json:"code"`
	Data []*TheatreInfo `json:"data"`
}

type TheatreInfo struct {
	CinemaName string `json:"cinemaName"`
	Address    string `json:"address"`
	CityName   string `json:"cityName"`
	CinemaId   uint   `json:"cinemaId"`
}

type Plus struct {
	Latitude  string
	Longitude string
	SessionId string
}

type CityList struct {
	Code uint        `json:"code"`
	Data *CityDetail `json:"data"`
}

type CityDetail struct {
	Hot []interface{} `json:"hot"`
	All map[string][]struct {
		Name string `json:"name"`
	} `json:"all"`
}

func (a Plus) CityList(ctx context.Context) []string {
	list := make([]string, 0)
	res, _ := ioutil.ReadFile("./city.json")
	fmt.Println(string(res))
	var c CityList
	err := json.Unmarshal(res, &c)
	if err != nil {
		fmt.Println(err)
		return list
	}
	for _, v := range c.Data.All {
		for _, j := range v {
			list = append(list, j.Name)
		}
	}
	return list
}

func (a Plus) MovieList(ctx context.Context) []uint {
	list := make([]uint, 0)
	return list
}

func (a Plus) FetchTheatre(ctx context.Context, targetUrl string, movieId uint64, cityName string) ([]*TheatreInfo, error) {
	var list []*TheatreInfo
	formData := map[string]string{
		"latitude":  a.Latitude,
		"longitude": a.Longitude,
		"city":      cityName,
		"movieId":   fmt.Sprint(movieId),
	}
	header := map[string]string{"sessionId": a.SessionId}
	nC := easy_http_client.NewHttpClient(targetUrl, easy_http_client.HttpClientMethodPost, header, 1, "")
	nC.SetPostFormData(ctx, formData)
	resp, err := nC.HttpClient(ctx)
	if err != nil {
		return list, err
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}
	var rep RespData
	err = json.Unmarshal(rest, &rep)
	return rep.Data, err
}

type TheatreInfoResp struct {
	Code uint               `json:"code"`
	Data *TheatreDetailInfo `json:"data"`
}

type TheatreDetailInfo struct {
	Cinema struct {
		CityName   string `json:"cityName"`
		CinemaAddr string `json:"cinemaAddr"`
		Longitude  string `json:"longitude"`
		Latitude   string `json:"latitude"`
	} `json:"cinema"`
	Features []struct {
		FeatureName string `json:"featureName"`
		FeatureDesc string `json:"featureDesc"`
	} `json:"features"`
}

func (a Plus) FetchTheatreInfo(ctx context.Context, targetUrl string, theatreId uint) (*TheatreDetailInfo, error) {
	var info TheatreDetailInfo
	formData := map[string]string{
		"cinemaId": fmt.Sprint(theatreId),
	}
	header := map[string]string{"sessionId": a.SessionId}
	nC := easy_http_client.NewHttpClient(targetUrl, easy_http_client.HttpClientMethodPost, header, 1, "")
	nC.SetPostFormData(ctx, formData)
	resp, err := nC.HttpClient(ctx)
	if err != nil {
		return &info, err
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &info, err
	}
	var rep TheatreInfoResp
	err = json.Unmarshal(rest, &rep)
	return rep.Data, err
}
