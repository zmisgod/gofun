package downloader

import (
	"errors"
	"github.com/zmisgod/goSpider/utils"
	"io"
	"log"
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
		log.Println(err)
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
	_, err = io.Copy(fd, resp.Body)
	return err
}

func (a *Downloader) SaveFile() error {
	if err := a.fetchHeader(a.Url); err != nil {
		return err
	}
	log.Println("start CreateFileReError")
	fd, err := utils.CreateFileReError(a.SavePath+a.SaveName)
	if err != nil {
		return err
	}
	log.Println("start fetchWriteBody")
	defer fd.Close()
	if err := a.fetchWriteBody(a.Url, fd); err != nil {
		return err
	}
	return nil
}
