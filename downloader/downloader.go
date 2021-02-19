package downloader

import (
	"context"
	"errors"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Downloader struct {
	Url                      string            `json:"url"`                         //下载的url
	SaveName                 string            `json:"save_name"`                   //保存文件名称
	SavePath                 string            `json:"save_path"`                   //保存的文件夹
	ProxyHost                string            `json:"proxy_host"`                  //设置http代理
	CustomHeader             map[string]string `json:"custom_header"`               //设置http的header
	Timeout                  int               `json:"timeout"`                     //设置超时时间
	BreakPointContinueUpload bool              `json:"break_point_continue_upload"` //是否需要支持断点续传
	DownloadRoutine          int               `json:"download_routine"`            //下载的协程
	fileSize                 int               `json:"file_size"`                   //文件的大小
	fd                       *os.File          `json:"fd"`                          //文件
}

type HttpMethod string

const (
	HTTPGet  HttpMethod = "GET"
	HTTPPost HttpMethod = "POST"
	HTTPHead HttpMethod = "HEAD"
)

const (
	HeaderRange                string = "Range"
	DefaultHeaderRangeTemplate string = "bytes=%d-%d"
	DefaultHeaderRangeStartId  int    = 0
	DefaultHeaderRangeEndId    int    = 4
)

func getDefaultHeaderRange() string {
	return getHeaderRange(DefaultHeaderRangeStartId, DefaultHeaderRangeEndId)
}

func getResponseHeaderRange() string {
	return fmt.Sprintf("bytes %d-%d", DefaultHeaderRangeStartId, DefaultHeaderRangeEndId)
}

func getHeaderRange(startId, endId int) string {
	return fmt.Sprintf(DefaultHeaderRangeTemplate, startId, endId)
}

const (
	DefaultHTTPTimeout     int = 10 //http超时时间
	DefaultDownloadRoutine int = 6  //下载的协程数量
)

var (
	ErrorUrlIsEmpty    = errors.New("url is empty")
	ErrorUrlIsNotFound = errors.New("url not found")
	ErrorFileIsError   = errors.New("file is error")
)

func NewDownloader(urlString string) (*Downloader, error) {
	if urlString == "" {
		return nil, ErrorUrlIsEmpty
	}
	saveName, err := genFileName(urlString)
	if err != nil {
		return nil, err
	}
	return &Downloader{
		Url:                      urlString,
		SavePath:                 "./",
		SaveName:                 saveName,
		Timeout:                  DefaultHTTPTimeout,
		DownloadRoutine:          DefaultDownloadRoutine,
		BreakPointContinueUpload: true,
		CustomHeader:             make(map[string]string, 0),
	}, nil
}

func genFileName(pathUrl string) (string, error) {
	saveName, err := url.Parse(pathUrl)
	if err != nil {
		return "", err
	}
	pathInfo := strings.Split(saveName.Path, "/")
	return pathInfo[len(pathInfo)-1], nil
}

func (a *Downloader) SetSaveName(name string) {
	if name != "" {
		a.SaveName = name
	}
}

func (a *Downloader) setFileSize(size int) {
	if size > 0 {
		a.fileSize = size
	}
}

func (a *Downloader) SetSavePath(path string) {
	if path != "" {
		a.SavePath = path
	}
}

func (a *Downloader) SetProxy(proxyHost string) {
	if proxyHost != "" {
		a.ProxyHost = proxyHost
	}
}

func (a *Downloader) SetTimeout(timeout int) {
	if timeout > 0 {
		a.Timeout = timeout
	}
}

func (a *Downloader) SupportBreakPointContinueUpload() {
	a.BreakPointContinueUpload = true
	a.DownloadRoutine = DefaultDownloadRoutine
}

func (a *Downloader) DisabledBreakPointContinueUpload() {
	a.BreakPointContinueUpload = false
	a.DownloadRoutine = 1
}

func (a *Downloader) SetCustomHeader(headers map[string]string) {
	if len(headers) > 0 {
		a.CustomHeader = headers
	} else {
		a.CustomHeader = make(map[string]string, 1)
		a.CustomHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36"
	}
}

//setDownloadFileInfo 获取下载文件的信息：文件大小、文件的真实名称
func (a *Downloader) setDownloadFileInfo(header http.Header) {
	cdData := header.Get("Content-Disposition")
	if cdData != "" {
		var re = regexp.MustCompile(`(?m)filename="(.*)"`)
		list := re.FindAllStringSubmatch(cdData, 10000)
		if len(list) > 0 && len(list[0]) >= 1 {
			a.SetSaveName(list[0][1])
		}
	}
	if a.BreakPointContinueUpload {
		crData := header.Get("Content-Range")
		if crData != "" {
			var re = regexp.MustCompile(fmt.Sprintf(`%s\/(.*)`, getResponseHeaderRange()))
			list := re.FindAllStringSubmatch(crData, 10000)
			if len(list) > 0 && len(list[0]) >= 1 {
				_n, _ := strconv.Atoi(list[0][1])
				a.setFileSize(_n)
			}
		}
	} else {
		clData := header.Get("Content-Length")
		if clData != "" {
			_n, _ := strconv.Atoi(clData)
			a.setFileSize(_n)
		}
	}
}

func (a *Downloader) SaveFile(ctx context.Context) error {
	a.initData()
	//获取文件真实名称、文件大小、是否支持断点续传
	if err := a.checkFileSupportBreakPointAndFileName(ctx, getDefaultHeaderRange()); err != nil {
		return err
	}
	if err := a.createFile(ctx); err != nil {
		return err
	}
	defer a.fd.Close()
	var wg sync.WaitGroup
	wg.Add(a.DownloadRoutine)
	per := a.fileSize / a.DownloadRoutine
	for i := 0; i < a.DownloadRoutine; i++ {
		startId := i * per
		endId := (i+1)*per - 1
		if i == (a.DownloadRoutine - 1) {
			endId = a.fileSize
		}
		go func(i, startId, endId int) {
			defer wg.Done()
			if err := a.doHttpRequest(ctx, startId, endId); err != nil {
				log.Printf("gId:%d range:%d-%d error:%v", i, startId, endId, err)
			} else {
				log.Printf("gId:%d range:%d-%d", i, startId, endId)
			}
		}(i, startId, endId)
	}
	wg.Wait()
	return nil
}

func (a *Downloader) initData() {
	a.SavePath = strings.TrimRight(a.SavePath, "/") + "/"
}

func (a *Downloader) prepareHTTPClient(context context.Context, targetURL string, method HttpMethod, rangeStr string) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(a.Timeout), //超时时间
	}
	if a.ProxyHost != "" {
		proxyStr, err := url.Parse(a.ProxyHost)
		if err != nil {
			return nil, nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyStr),
		}
	}
	request, err := http.NewRequest(string(method), targetURL, nil)
	if err != nil {
		return nil, nil, err
	}
	if len(a.CustomHeader) > 0 {
		for k, v := range a.CustomHeader {
			request.Header.Set(k, v)
		}
	}
	if rangeStr != "" {
		request.Header.Set(HeaderRange, rangeStr)
	}
	return client, request, nil
}

func (a *Downloader) checkFileSupportBreakPointAndFileName(ctx context.Context, rangeStr string) error {
	httpClient, httpRequest, err := a.prepareHTTPClient(ctx, a.Url, HTTPGet, rangeStr)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//检查文件是否支持断点续传
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return ErrorUrlIsNotFound
	}
	if resp.StatusCode != http.StatusPartialContent {
		a.DisabledBreakPointContinueUpload()
	}
	//获取文件的真实文件名称
	a.setDownloadFileInfo(resp.Header)
	return nil
}

func (a *Downloader) doHttpRequest(ctx context.Context, startId, endId int) error {
	rangeStr := getHeaderRange(startId, endId)
	httpClient, httpRequest, err := a.prepareHTTPClient(ctx, a.Url, HTTPGet, rangeStr)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if a.fd == nil {
		return ErrorFileIsError
	}
	fd := a.fd
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = fd.WriteAt(data, int64(startId))
	return err
}

func (a *Downloader) createFile(ctx context.Context) error {
	fd, err := utils.CreateFileReError(a.SavePath + a.SaveName)
	if err != nil {
		if err == utils.ErrorFileExists {
			if !a.BreakPointContinueUpload {
				return err
			}
			fd, _ = os.OpenFile(a.SavePath+a.SaveName, os.O_RDWR, 0666)
			fileStatus, _ := fd.Stat()
			_, err := fd.Seek(fileStatus.Size(), 0)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	a.fd = fd
	return nil
}
