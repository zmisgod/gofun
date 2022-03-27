package easy_http_client

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClientMethod string

const (
	HttpClientMethodGet  = "GET"
	HttpClientMethodPost = "POST"
	HttpClientMethodHead = "HEAD"
)

type HttpClientPostType uint8

const (
	HttpClientPostTypeForm = 1
	HttpClientPostTypeJson = 2
)

func NewHttpClient(targetUrl string, method HttpClientMethod, header map[string]string, timeOut int, proxy string) *HttpClient {
	if header == nil {
		header = make(map[string]string, 0)
	}
	return &HttpClient{
		TargetUrl: targetUrl,
		Method:    method,
		Header:    header,
		Timeout:   timeOut,
		Proxy:     proxy,
	}
}

type HttpClient struct {
	TargetUrl string
	Method    HttpClientMethod
	Header    map[string]string
	Timeout   int
	Proxy     string
	Body      io.Reader
}

func (a *HttpClient) SetPostFormData(ctx context.Context, body map[string]string) {
	data := url.Values{}
	for k, v := range body {
		data.Set(k, v)
	}
	a.Body = strings.NewReader(data.Encode())
	a.Header["Content-Type"] = "application/x-www-form-urlencoded"
}

func (a *HttpClient) SetJsonData(ctx context.Context, jsonString string) {
	a.Body = strings.NewReader(jsonString)
	a.Header["Content-Type"] = "application/json"
}

func (a *HttpClient) HttpClient(ctx context.Context) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(a.Timeout), //超时时间
	}
	transPort := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if a.Proxy != "" {
		proxyStr, err := url.Parse(a.Proxy)
		if err != nil {
			return nil, err
		}
		transPort.Proxy = http.ProxyURL(proxyStr)
	}
	client.Transport = transPort
	request, err := http.NewRequest(string(a.Method), a.TargetUrl, a.Body)
	if err != nil {
		return nil, err
	}
	if len(a.Header) > 0 {
		for k, v := range a.Header {
			request.Header.Set(k, v)
		}
	}
	request = request.WithContext(ctx)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
