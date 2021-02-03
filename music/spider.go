package music

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	FetchMusicTypeQQ     = 1 //qq 音乐
	FetchMusicTypeNameQQ = "QQ音乐"
	FetchMusicTypeWY     = 2 //网易云
	FetchMusicTypeNameWY = "网易云音乐"
	QQMusicGetSignUrl = "http://127.0.0.1:3000/getSign"
)

type FetchMusicSpider struct {
	OriginalURL string

	fetchType uint8  //需要爬取的类型
	fetchURL  string //需要爬取的URL

	musicURL string //结果的url

	musicID     string
	musicName   string
	musicSinger string
}

func NewFetchMusic(url string) (*FetchMusicSpider, error) {
	if url == "" {
		return nil, errors.New("url为空")
	}
	return &FetchMusicSpider{
		OriginalURL: url,
	}, nil
}

func (a *FetchMusicSpider) GetDownloadURL(ctx context.Context) (string, error) {
	if err := a.parseLinkURL(); err != nil {
		return "", err
	}
	if err := a.fetchFromURL(ctx); err != nil {
		return "", err
	}
	if a.musicURL == "" {
		return "", errors.New("获取失败")
	}
	return a.musicURL, nil
}

func (a *FetchMusicSpider) fetchFromURL(ctx context.Context) error {
	if a.fetchType == FetchMusicTypeQQ {
		if err := a.fetchQQMusic(ctx); err != nil {
			return err
		}
	} else if a.fetchType == FetchMusicTypeWY {
		if err := a.fetchWYMusic(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *FetchMusicSpider) parseLinkURL() error {
	//平台
	var platFormReg = regexp.MustCompile(`@[\x{4e00}-\x{9fa5}A-Za-z0-9_]+`)
	platFormRes := platFormReg.FindStringSubmatch(a.OriginalURL)
	if len(platFormRes) == 0 {
		return errors.New("无法解析")
	}
	platform := strings.Replace(platFormRes[0], "@", "", 100)
	if platform == FetchMusicTypeNameQQ {
		a.fetchType = FetchMusicTypeQQ
	} else if platform == FetchMusicTypeNameWY {
		a.fetchType = FetchMusicTypeWY
	} else {
		return errors.New("不支持的平台")
	}
	if a.fetchType == FetchMusicTypeQQ {
		authorReg := regexp.MustCompile(`(?m)^[\x{4e00}-\x{9fa5}A-Za-z0-9_]+`)
		authorRes := authorReg.FindStringSubmatch(a.OriginalURL)
		if len(authorRes) > 0 {
			a.musicSinger = authorRes[0]
		}
	} else if a.fetchType == FetchMusicTypeWY {
		authorReg := regexp.MustCompile(`分享(.*)的`)
		authorRes := authorReg.FindStringSubmatch(a.OriginalURL)
		if len(authorRes) > 0 {
			a.musicSinger = authorRes[1]
		}
	}
	//歌曲名
	songNameReg := regexp.MustCompile(`(?m)《(.*)》`)
	songNameRes := songNameReg.FindStringSubmatch(a.OriginalURL)
	if len(songNameRes) > 0 {
		a.musicName = songNameRes[1]
	}
	//url
	urlReg := regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	urlRes := urlReg.FindStringSubmatch(a.OriginalURL)
	if len(urlRes) > 0 && urlRes[0] != "" {
		a.fetchURL = urlRes[0]
	} else {
		return errors.New("获取失败")
	}
	return nil
}

func (a *FetchMusicSpider) getQQMusicID() error {
	client := &http.Client{
		Timeout: time.Second * 10, //超时时间
	}
	request, err := http.NewRequest(http.MethodGet, a.fetchURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36")
	request.Header.Set("origin", "https://y.qq.com")
	request.Header.Set("referer", "https://y.qq.com/portal/player.html")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	restBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	urlReg := regexp.MustCompile(`http-equiv="refresh" content="0; url=(.*)"`)
	urls := urlReg.FindStringSubmatch(string(restBody))
	if len(urls) != 2 {
		return errors.New("获取QQ音乐链接失败")
	}
	linkUrl := html.UnescapeString(urls[1])
	midReg := regexp.MustCompile(`mid=(.*)&`)
	midRes := midReg.FindStringSubmatch(linkUrl)
	if len(midRes) != 2 {
		return errors.New("获取QQ音乐ID失败")
	}
	a.musicID = midRes[1]
	return nil
}

func (a *FetchMusicSpider) fetchQQMusic(ctx context.Context) error {
	err := a.getQQMusicID()
	if err != nil {
		return err
	}
	sign := a.getQQMusicRequestSign(ctx)
	if sign == "" {
		return errors.New("sign error")
	}
	link := a.getQQMusicUrl()
	if link == "" {
		return errors.New("link is empty")
	}
	a.musicURL = link
	return nil
}

func (a *FetchMusicSpider) getQQMusicRequestSign(ctx context.Context) string {
	req := fmt.Sprintf("{\"req\":{\"module\":\"CDN.SrfCdnDispatchServer\",\"method\":\"GetCdnDispatch\",\"param\":{\"guid\":\"52278914\",\"calltype\":0,\"userip\":\"\"}},\"req_0\":{\"module\":\"vkey.GetVkeyServer\",\"method\":\"CgiGetVkey\",\"param\":{\"guid\":\"52278914\",\"songmid\":[\"%s\"],\"songtype\":[0],\"uin\":\"0\",\"loginflag\":0,\"platform\":\"20\"}},\"comm\":{\"uin\":0,\"format\":\"json\",\"ct\":24,\"cv\":0}}", a.musicID)
	sign := a.FetchQQMusicSign(ctx, a.fetchType, req)
	if sign == "" {
		return ""
	}
	return sign
}

func (a *FetchMusicSpider) FetchQQMusicSign(ctx context.Context, cType uint8, body string) string {
	resp, err := http.Post(QQMusicGetSignUrl, "application/json", strings.NewReader(body))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	type NodeResult struct {
		Code int `json:"code"`
		Data string `json:"data"`
	}
	var nodeResult NodeResult
	_  = json.Unmarshal(rest, &nodeResult)
	if nodeResult.Code == 200 {
		return nodeResult.Data
	}
	return ""
}

func (a *FetchMusicSpider) getQQMusicUrl() string {
	return ""
}

func (a *FetchMusicSpider) fetchWYMusic(ctx context.Context) error {
	IDReg := regexp.MustCompile(`id=(.*)&userid`)
	res := IDReg.FindStringSubmatch(a.fetchURL)
	if len(res) > 0 {
		a.musicID = res[1]
	} else {
		return errors.New("音乐ID获取失败")
	}
	a.musicURL = fmt.Sprintf("http://music.163.com/song/media/outer/url?id=%s.mp3", a.musicID)
	return nil
}
