package downloader

import (
	"bufio"
	"errors"
	"github.com/zmisgod/gofun/utils"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Downloader struct {
	Url      string `json:"url"`
	SaveName string `json:"save_name"`
	SavePath string `json:"save_path"`
}

var (
	ErrorUrlIsEmpty = errors.New("url is empty")
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
		Url:      urlString,
		SavePath: "./",
		SaveName: saveName,
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

func (a *Downloader) SetSavePath(path string) {
	if path != "" {
		a.SavePath = path
	}
}

func (a *Downloader) fetchHeader(urlStr string) error {
	resp, err := http.Head(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resData := resp.Header.Get("Content-Disposition")
	if resData != "" {
		var re = regexp.MustCompile(`(?m)filename="(.*)"`)
		list := re.FindAllStringSubmatch(resData, 10000)
		if len(list) > 0 && len(list[0]) >= 1 {
			a.SetSaveName(list[0][1])
		}
	}
	return nil
}

func (a *Downloader) fetchWriteBody(urlStr string, fd *os.File) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	wt := bufio.NewWriter(fd)
	_, err = io.Copy(wt, resp.Body)
	_ = wt.Flush()
	return err
}

func (a *Downloader) SaveFile() error {
	a.initData()
	if err := a.fetchHeader(a.Url); err != nil {
		return err
	}
	fd, err := utils.CreateFileReError(a.SavePath + a.SaveName)
	if err != nil {
		return err
	}
	defer fd.Close()
	if err := a.fetchWriteBody(a.Url, fd); err != nil {
		return err
	}
	return nil
}

func (a *Downloader) initData() {
	a.SavePath = strings.TrimRight(a.SavePath, "/") + "/"
}
