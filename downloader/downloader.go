package downloader

import (
	"context"
	"crypto/tls"
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
	Url      string            `json:"url"`       //下载的url
	fileSize int               `json:"file_size"` //文件的大小
	fd       *os.File          `json:"fd"`        //文件指针
	option   options           `json:"option"`    //参数
	fileMap  map[uint]*os.File `json:"file_map"`  //文件指针map for 断点续传
}

type options struct {
	SaveName        string            `json:"save_name"`        //保存文件名称
	SavePath        string            `json:"save_path"`        //保存的文件夹
	ProxyHost       string            `json:"proxy_host"`       //设置http代理
	CustomHeader    map[string]string `json:"custom_header"`    //设置http的header
	Timeout         int               `json:"timeout"`          //设置超时时间
	DownloadRoutine int               `json:"download_routine"` //下载的协程
	BreakPoint      bool              `json:"break_point"`      //是否需要支持断点续传
	TryTimes        int               `json:"try_times"`        //失败重试次数
	StrategyWait    bool              `json:"strategy_wait"`    //策略等待
}

type Option func(*options)

type HttpMethod string

const (
	HTTPGet  HttpMethod = "GET"
	HTTPPost HttpMethod = "POST"
	HTTPHead HttpMethod = "HEAD"
)

const (
	HeaderRange                string = "Range"
	DefaultHeaderRangeTemplate string = "bytes=%d-%d"
	DefaultHeaderRangeStartId  int    = 3
	DefaultHeaderRangeEndId    int    = 4
	DefaultTryTimes            int    = 30
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
	DefaultHTTPTimeout             int = 10 //http超时时间
	DefaultDownloadRoutine         int = 6  //下载的协程数量
	DefaultDisabledDownloadRoutine int = 1  //不支持下载的协程数量
)

var (
	ErrorUrlIsEmpty    = errors.New("url is empty")
	ErrorUrlIsNotFound = errors.New("url not found")
	ErrorFileIsError   = errors.New("file is error")
)

func SetTryTimes(times int ) Option {
	return func(o *options) {
		o.TryTimes = times
	}
}

func SetStrategyWait(isSmart bool) Option {
	return func(o *options) {
		o.StrategyWait = isSmart
	}
}

func SetSaveFileName(name string) Option {
	return func(o *options) {
		o.SaveName = name
	}
}

func SetSavePath(path string) Option {
	return func(o *options) {
		o.SavePath = path
	}
}

func SetProxyHost(proxyHost string) Option {
	return func(o *options) {
		o.ProxyHost = proxyHost
	}
}

func SetCustomHeader(header map[string]string) Option {
	return func(o *options) {
		o.CustomHeader = header
	}
}

func SetTimeout(timeout int) Option {
	return func(o *options) {
		o.Timeout = timeout
	}
}

func SetBreakPoint(isNeed bool) Option {
	return func(o *options) {
		o.BreakPoint = isNeed
	}
}

func SetDownloadRoutine(num int) Option {
	return func(o *options) {
		o.DownloadRoutine = num
	}
}

func NewDownloader(urlString string, option ...Option) (*Downloader, error) {
	if urlString == "" {
		return nil, ErrorUrlIsEmpty
	}
	var op options
	if len(option) > 0 {
		for _, opt := range option {
			opt(&op)
		}
	}
	if op.SaveName == "" {
		_name, err := genFileName(urlString)
		if err != nil {
			return nil, err
		}
		op.SaveName = _name
	}
	if op.SavePath == "" {
		op.SavePath = "./"
	}
	if len(op.CustomHeader) == 0 {
		op.CustomHeader = make(map[string]string, 0)
		op.CustomHeader[utils.UserAgentName] = utils.UserAgentString
	}
	if op.DownloadRoutine == 0 {
		op.DownloadRoutine = DefaultDownloadRoutine
	}
	if op.Timeout == 0 {
		op.Timeout = DefaultHTTPTimeout
	}
	if op.TryTimes == 0 {
		op.TryTimes = DefaultTryTimes
	}
	return &Downloader{
		Url:    urlString,
		option: op,
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

func (a *Downloader) setSaveName(name string) {
	if name != "" {
		a.option.SaveName = name
	}
}

func (a *Downloader) setFileSize(size int) {
	if size > 0 {
		a.fileSize = size
	}
}

func (a *Downloader) setSavePath(path string) {
	if path != "" {
		a.option.SavePath = path
	}
}

func (a *Downloader) setProxy(proxyHost string) {
	if proxyHost != "" {
		a.option.ProxyHost = proxyHost
	}
}

func (a *Downloader) setTimeout(timeout int) {
	if timeout > 0 {
		a.option.Timeout = timeout
	}
}

func (a *Downloader) setDownloadRoutine(num int) {
	if num > 0 {
		a.option.DownloadRoutine = DefaultDownloadRoutine
	}
}

func (a *Downloader) setCustomHeader(headers map[string]string) {
	if len(headers) > 0 {
		a.option.CustomHeader = headers
	} else {
		a.option.CustomHeader = make(map[string]string, 1)
		a.option.CustomHeader["User-Agent"] = utils.UserAgentString
	}
}

//setDownloadFileInfo 获取下载文件的信息：文件大小、文件的真实名称
func (a *Downloader) setDownloadFileInfo(header http.Header) {
	cdData := header.Get("Content-Disposition")
	if cdData != "" {
		var re = regexp.MustCompile(`(?m)filename="(.*)"`)
		list := re.FindAllStringSubmatch(cdData, 10000)
		if len(list) > 0 && len(list[0]) >= 1 {
			a.setSaveName(list[0][1])
		}
	}
	crData := header.Get("Content-Range")
	if crData != "" {
		var re = regexp.MustCompile(fmt.Sprintf(`%s\/(.*)`, getResponseHeaderRange()))
		list := re.FindAllStringSubmatch(crData, 10000)
		if len(list) > 0 && len(list[0]) >= 1 {
			_n, _ := strconv.Atoi(list[0][1])
			a.setFileSize(_n)
		}
	}
	if a.fileSize == 0 {
		clData := header.Get("Content-Length")
		if clData != "" {
			_n, _ := strconv.Atoi(clData)
			a.setFileSize(_n)
		}
	}
}

func (a *Downloader) SaveFile(ctx context.Context) error {
	a.initData()
	//获取文件真实名称、文件大小、是否支持多协程下载
	if err := a.checkFileSupportMultiRoutineAndFileName(ctx, getDefaultHeaderRange()); err != nil {
		return err
	}
	if a.option.BreakPoint {
		return a.doMultiDownloadWithBreakPoint(ctx)
	} else {
		return a.doMultiDownloadWithoutBreakPoint(ctx)
	}
}

func (a *Downloader) initData() {
	a.option.SavePath = strings.TrimRight(a.option.SavePath, "/") + "/"
}

func (a *Downloader) prepareHTTPClient(ctx context.Context, targetURL string, method HttpMethod, rangeStr string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(a.option.Timeout), //超时时间
	}
	transPort := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if a.option.ProxyHost != "" {
		proxyStr, err := url.Parse(a.option.ProxyHost)
		if err != nil {
			return nil, err
		}
		transPort.Proxy = http.ProxyURL(proxyStr)
	}
	client.Transport = transPort
	request, err := http.NewRequest(string(method), targetURL, nil)
	if err != nil {
		return nil, err
	}
	if len(a.option.CustomHeader) > 0 {
		for k, v := range a.option.CustomHeader {
			request.Header.Set(k, v)
		}
	}
	if rangeStr != "" {
		request.Header.Set(HeaderRange, rangeStr)
	}
	request = request.WithContext(ctx)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *Downloader) checkFileSupportMultiRoutineAndFileName(ctx context.Context, rangeStr string) error {
	resp, err := a.prepareHTTPClient(ctx, a.Url, HTTPGet, rangeStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//检查文件是否支持range
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		log.Println(fmt.Sprintf("http response statusCode is :%d", resp.StatusCode))
		return ErrorUrlIsNotFound
	}
	if resp.StatusCode != http.StatusPartialContent {
		a.setDownloadRoutine(1)
	}
	//获取文件的真实文件名称
	a.setDownloadFileInfo(resp.Header)
	return nil
}

func (a *Downloader) doHttpRequest(ctx context.Context, startId, endId int) error {
	rangeStr := getHeaderRange(startId, endId)
	resp, err := a.prepareHTTPClient(ctx, a.Url, HTTPGet, rangeStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if a.fd == nil {
		return ErrorFileIsError
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.New(fmt.Sprintf("startId %d- endId %d response status is not valid,now is %d", startId, endId, resp.StatusCode))
	}
	result, _ := ioutil.ReadAll(resp.Body)
	_, err = a.fd.WriteAt(result, int64(startId))
	return err
}

func getPerFileName(name string, i uint64) string {
	res := strings.Split(name, ".")
	fmt.Println(res)
	return ""
}

func (a *Downloader) createFile() error {
	fd, err := utils.CreateFileReError(a.option.SavePath + a.option.SaveName)
	if err != nil {
		if err == utils.ErrorFileExists {
			if a.option.DownloadRoutine > DefaultDisabledDownloadRoutine {
				return err
			}
			fd, _ = os.OpenFile(a.option.SavePath+a.option.SaveName, os.O_RDWR, 0666)
		} else {
			return err
		}
	} else {
		if err := fd.Truncate(int64(a.fileSize)); err != nil {
			return err
		}
	}
	a.fd = fd
	return nil
}

func (a *Downloader) createFileByPer(i uint64, size uint64) (int64, error) {
	var fileSize int64
	fd, err := utils.CreateFileReError(a.option.SavePath + a.option.SaveName + fmt.Sprintf(".tmp%d", i))
	if err != nil {
		if err == utils.ErrorFileExists {
			fd, _ = os.OpenFile(a.option.SavePath+a.option.SaveName, os.O_RDWR, 0666)
			fileStatus, _ := fd.Stat()
			fileSize = fileStatus.Size()
			_, err := fd.Seek(fileSize, 0)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	a.fd = fd
	return fileSize, nil
}

func (a *Downloader) doMultiDownloadWithoutBreakPoint(ctx context.Context) error {
	if err := a.createFile(); err != nil {
		return err
	}
	defer a.fd.Close()
	childCtx, _ := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(a.option.DownloadRoutine)
	if a.option.StrategyWait {
		time.Sleep(3 * time.Second)
	}
	per := a.fileSize / a.option.DownloadRoutine
	for i := 0; i < a.option.DownloadRoutine; i++ {
		startId := i * per
		endId := (i+1)*per - 1
		if i == (a.option.DownloadRoutine - 1) {
			endId = a.fileSize
		}
		go func(ctx context.Context, i, startId, endId int) {
			defer wg.Done()
			for j := 0; j < a.option.TryTimes; j++ {
				if a.option.StrategyWait {
					time.Sleep(time.Duration(utils.Rand(3, 10)) * time.Second)
				}
				if err := a.doHttpRequest(ctx, startId, endId); err != nil {
					log.Printf("gId:%d range:%d-%d error:%v try:%d", i, startId, endId, err, j)
				} else {
					log.Printf("\033[32m gId:%d range:%d-%d \033[0m", i, startId, endId)
					break
				}
			}
			if a.option.StrategyWait {
				time.Sleep(time.Duration(i)*time.Second)
			}
		}(childCtx, i, startId, endId)
	}
	wg.Wait()
	return nil
}

//TODO
func (a *Downloader) doMultiDownloadWithBreakPoint(ctx context.Context) error {
	//per := a.fileSize / a.option.DownloadRoutine
	//fmt.Println(per)
	//for i := 0; i < a.option.DownloadRoutine; i++ {
	//
	//}
	return errors.New("not impl")
}
