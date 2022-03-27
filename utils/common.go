package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateFile(textName string) {
	_, err := os.Stat(textName)
	if err != nil {
		file, err := os.Create(textName)
		CheckError(err)
		defer file.Close()
	}
}

var ErrorFileExists error = errors.New("file exits")

//CreateFileReError 创建文件
func CreateFileReError(textName string) (*os.File, error) {
	_, err := os.Stat(textName)
	if err != nil {
		file, err := os.Create(textName)
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		return nil, ErrorFileExists
	}
}

//CreateFolder 创建文件夹
func CreateFolder(folderName string) (bool, error) {
	checkFolderNotExists := CheckPathIsNotExists(folderName)
	if checkFolderNotExists {
		err := os.MkdirAll(folderName, 0777)
		if err != nil {
			return false, err
		}
		log.Printf("create floder %s successful\n", folderName)
		return true, nil
	}
	return false, nil
}

//CheckPathIsNotExists 检查文件是否存在 返回true 不存在， false 存在
func CheckPathIsNotExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type HttpClientMethod string

const (
	HttpClientMethodGet  HttpClientMethod = "GET"
	HttpClientMethodPost HttpClientMethod = "POST"
)

var DefaultUserAgent map[string]string = map[string]string{
	UserAgentName: UserAgentString,
}

func Rand(min, max int) int {
	if min > max {
		return max
	}
	if int31 := 1<<31 - 1; max > int31 {
		return min
	}
	if min == max {
		return min
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

const UserAgentName = "User-Agent"
const UserAgentString = "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"

var ContentTypeFormUrlEncode = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
}

var ContentTypeJson = map[string]string{
	"Content-Type": "application/json",
}

func HttpClient(ctx context.Context, method HttpClientMethod, targetURL string, timeout time.Duration, proxy string, body io.ReadCloser, customHeader map[string]string) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout, //超时时间
	}
	transPort := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, timeout) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = c.SetDeadline(time.Now().Add(timeout)) //设置发送接收数据超时
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	if proxy != "" {
		proxyStr, err := url.Parse(proxy)
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
	if len(customHeader) > 0 {
		for k, v := range customHeader {
			request.Header.Set(k, v)
		}
	}
	if body != nil {
		request.Body = body
	}
	request = request.WithContext(ctx)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

const (
	DefaultHeaderRangeTemplate string = "bytes=%d-%d"
	HeaderRange                string = "Range"
)

func GetHeaderRange(startId, endId int) string {
	return fmt.Sprintf(DefaultHeaderRangeTemplate, startId, endId)
}
