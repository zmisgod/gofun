package music

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

const (
	FetchMusicTypeQQ     = 1 //qq 音乐
	FetchMusicTypeNameQQ = "QQ音乐"
	FetchMusicTypeWY     = 2 //网易云
	FetchMusicTypeNameWY = "网易云音乐"
	QQMusicGetSignUrl    = "http://127.0.0.1:20050/getSign"
)

var (
	ErrorLinkIsEmpty                = errors.New("link is empty")
	ErrorUrlIsEmpty                 = errors.New("url为空")
	ErrorCanNotParseYourLink        = errors.New("can not parse your input link")
	ErrorCanNotSupportMusicPlatform = errors.New("can not support music platform")
)

type FetchMusic interface {
	Fetch(ctx context.Context, musicInfo *LinkInfo) (string, error)
}

//LinkInfo 通过分享链接获取的音乐数据
type LinkInfo struct {
	Name   string `json:"name"`
	Singer string `json:"singer"`
	Url    string `json:"url"`
}

type FetchMusicSpider struct {
	OriginalURL string
	fetchType   uint8  //需要爬取的类型
	musicURL    string //结果的url
	LinkInfo    *LinkInfo
}

func NewFetchMusic(url string) (*FetchMusicSpider, error) {
	if url == "" {
		return nil, ErrorUrlIsEmpty
	}
	return &FetchMusicSpider{
		OriginalURL: url,
	}, nil
}

func (a *FetchMusicSpider) GetDownloadURL(ctx context.Context) (string, error) {
	a.LinkInfo = new(LinkInfo)
	if err := a.parseLinkURL(); err != nil {
		return "", err
	}
	if a.LinkInfo.Url == "" || a.LinkInfo.Name == "" || a.LinkInfo.Singer == "" {
		return "", errors.New("部分参数缺失")
	}
	if err := a.fetchFromURL(ctx); err != nil {
		return "", err
	}
	return a.musicURL, nil
}

func (a *FetchMusicSpider) fetchFromURL(ctx context.Context) error {
	var r FetchMusic
	if a.fetchType == FetchMusicTypeQQ {
		r = new(QQMusic)
	} else if a.fetchType == FetchMusicTypeWY {
		r = new(NetEastMusic)
	}
	if r == nil {
		return errors.New("不支持的类型")
	}
	obj, err := r.Fetch(ctx, a.LinkInfo)
	if err != nil {
		return err
	}
	if obj == "" {
		return errors.New("获取失败")
	}
	a.musicURL = obj
	return nil
}

func (a *FetchMusicSpider) parseLinkURL() error {
	//平台
	var platFormReg = regexp.MustCompile(`@[\x{4e00}-\x{9fa5}A-Za-z0-9_]+`)
	platFormRes := platFormReg.FindStringSubmatch(a.OriginalURL)
	if len(platFormRes) == 0 {
		return ErrorCanNotParseYourLink
	}
	platform := strings.Replace(platFormRes[0], "@", "", 100)
	if platform == FetchMusicTypeNameQQ {
		a.fetchType = FetchMusicTypeQQ
	} else if platform == FetchMusicTypeNameWY {
		a.fetchType = FetchMusicTypeWY
	} else {
		return ErrorCanNotSupportMusicPlatform
	}
	//歌手
	if a.fetchType == FetchMusicTypeQQ {
		authorReg := regexp.MustCompile(`(?m)^[\x{4e00}-\x{9fa5}A-Za-z0-9_]+`)
		authorRes := authorReg.FindStringSubmatch(a.OriginalURL)
		if len(authorRes) > 0 {
			a.LinkInfo.Singer = authorRes[0]
		}
	} else if a.fetchType == FetchMusicTypeWY {
		authorReg := regexp.MustCompile(`分享(.*)的`)
		authorRes := authorReg.FindStringSubmatch(a.OriginalURL)
		if len(authorRes) > 0 {
			a.LinkInfo.Singer = authorRes[1]
		}
	}
	//歌曲名
	songNameReg := regexp.MustCompile(`(?m)《(.*)》`)
	songNameRes := songNameReg.FindStringSubmatch(a.OriginalURL)
	if len(songNameRes) > 0 {
		a.LinkInfo.Name = songNameRes[1]
	}
	//url
	urlReg := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	urlRes := urlReg.FindStringSubmatch(a.OriginalURL)
	if len(urlRes) > 0 && urlRes[0] != "" {
		a.LinkInfo.Url = urlRes[0]
	}
	return nil
}
