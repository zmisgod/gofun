package downloader

import (
	"errors"
	"net/url"
	"strings"
)

type Downloader struct {
	Url string `json:"url"`
	SaveName string `json:"save_name"`
	SavePath string `json:"save_path"`
}

var (
	ErrorUrlIsEmpty = errors.New("url is empty")
)

func NewDownloader(urlString string) (*Downloader, error){
	if urlString == "" {
		return nil, ErrorUrlIsEmpty
	}
	saveName, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	pathInfo := strings.Split(saveName.Path, "/")
	return &Downloader{
		Url: urlString,
		SavePath: "./",
		SaveName: pathInfo[len(pathInfo)-1],
	}, nil
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

func (a *Downloader) SaveFile() error {

}