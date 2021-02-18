package music

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	FetchMusicTypeQQ     = 1 //qq 音乐
	FetchMusicTypeNameQQ = "QQ音乐"
	FetchMusicTypeWY     = 2 //网易云
	FetchMusicTypeNameWY = "网易云音乐"
	QQMusicGetSignUrl    = "http://127.0.0.1:20050/getSign"
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
	//err := a.getQQMusicID()
	//if err != nil {
	//	return err
	//}
	//sign := a.getQQMusicRequestSign(ctx)
	//if sign == "" {
	//	return errors.New("sign error")
	//}
	link := a.getQQMusicUrl(ctx)
	if link == "" {
		return errors.New("link is empty")
	}
	a.musicURL = link
	fmt.Println(link)
	fmt.Println("ok-----")
	return nil
}

func (a *FetchMusicSpider) getQQMusicRequestSign(ctx context.Context) (string, string, int) {
	qq := 838881690
	//req := fmt.Sprintf("{\"req\":{\"module\":\"CDN.SrfCdnDispatchServer\",\"method\":\"GetCdnDispatch\",\"param\":{\"guid\":\"6509348615\",\"calltype\":0,\"userip\":\"\"}},\"req_0\":{\"module\":\"vkey.GetVkeyServer\",\"method\":\"CgiGetVkey\",\"param\":{\"guid\":\"6509348615\",\"songmid\":[\"%s\"],\"songtype\":[0],\"uin\":\"%d\",\"loginflag\":%d,\"platform\":\"20\"}},\"comm\":{\"uin\":%d,\"format\":\"json\",\"ct\":24,\"cv\":0}}", a.musicID, qq, 1, qq)
	req := `{"req":{"module":"CDN.SrfCdnDispatchServer","method":"GetCdnDispatch","param":{"guid":"6509348615","calltype":0,"userip":""}},"req_0":{"module":"vkey.GetVkeyServer","method":"CgiGetVkey","param":{"guid":"6509348615","songmid":["0004Fq5m1or8Sq"],"songtype":[0],"uin":"838881690","loginflag":1,"platform":"20"}},"comm":{"uin":838881690,"format":"json","ct":24,"cv":0}}`
	sign := a.FetchQQMusicSign(ctx, a.fetchType, req)
	if sign == "" {
		return "", "", 0
	}
	return sign, req, qq
}

func (a *FetchMusicSpider) FetchQQMusicSign(ctx context.Context, cType uint8, body string) string {
	postData := make(url.Values)
	postData["needSign"] = []string{body}
	resp, err := http.PostForm(QQMusicGetSignUrl,  postData)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	type NodeResult struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	var nodeResult NodeResult
	_ = json.Unmarshal(rest, &nodeResult)
	if nodeResult.Code == 200 {
		return nodeResult.Data
	}
	return ""
}

func parseUrl(sText string) string {
	sText = strings.Replace(sText, "&#58;", ":", 10000)
	sText = strings.Replace(sText, "&#47;", "/", 10000)
	sText = strings.Replace(sText, "&#46;", ".", 10000)
	sText = strings.Replace(sText, "&#35;", "#", 10000)
	sText = strings.Replace(sText, "&#61;", "=", 10000)
	sText = strings.Replace(sText, "&#38;", "&", 10000)
	sText = strings.Replace(sText, "&#34;", "\"", 10000)
	return sText
}

func getMusicMID(urlStr string) string {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	resultUrl := ""
	doc.Find("head meta").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		str, exists := selection.Attr("content")
		if exists {
			var re = regexp.MustCompile(`url=(.*)`)
			list := re.FindAllStringSubmatch(str, 1000)
			if len(list) > 0 {
				if len(list[0]) >= 2 {
					resultUrl = list[0][1]
				}
			}
			return true
		}
		return true
	})
	return resultUrl
}

type QQMusicObj struct {
	SongId       uint   `json:"songid"`
	SongMId      string `json:"songmid"`
	SongName     string `json:"songname"`
	SongTitle    string `json:"songtitle"`
	AlbumId      uint   `json:"albumid"`
	AlbumMId     string `json:"albummid"`
	AlbumName    string `json:"albumname"`
	AlbumPMId    string `json:"albumpmid"`
	Mid          string `json:"mid"`
	StrMediaMid  string `json:"strMediaMid"`
	SongSubTitle string `json:"songsubtitle"`
	SingerId     uint   `json:"singerid"`
	SingerMId    string `json:"singermid"`
}

func getMusicInfo(urlStr string) *QQMusicObj {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	var resultObj QQMusicObj
	doc.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		str, _ := selection.Html()
		if str != "" {
			str = parseUrl(str)
			var re = regexp.MustCompile(`var g_SongData = (.*);`)
			list := re.FindAllStringSubmatch(str, 1000)
			if len(list) > 0 {
				if len(list[0]) >= 2 {
					_ = json.Unmarshal([]byte(list[0][1]), &resultObj)
				}
			}
		}
		return true
	})
	return &resultObj
}

type MusicsFcg struct {
	Code uint `json:"code"`
	Req0 struct {
		Code uint `json:"code"`
		Data struct {
			Midurlinfo []struct{
				PUrl string `json:"purl"`
			} `json:"midurlinfo"`
			Sip          []string `json:"sip"`
		} `json:"data"`
	} `json:"req_0"`
}

func (a *FetchMusicSpider) getQQMusicUrl(ctx context.Context) string {
	//songId :=  getMusicMID(a.fetchURL)
	songId := "0004Fq5m1or8Sq"
	musicObj := getMusicInfo(fmt.Sprintf("https://y.qq.com/n/yqq/song/%s.html?ADTAG=h5_playsong&no_redirect=1", songId))
	if musicObj != nil {
		a.musicID = musicObj.SongMId
		getSign, dataStr, qq := a.getQQMusicRequestSign(ctx)
		rest :="https://u.y.qq.com/cgi-bin/musics.fcg?-=getplaysongvkey3047664478189791&g_tk=743553847&sign="+getSign+"&loginUin="+strconv.Itoa(qq)+"&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&data="+ url.PathEscape(dataStr)
		fmt.Println(rest)
		resp, err := http.Get(rest)
		if err != nil {
			log.Println(err)
			return ""
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return ""
		}
		var musicFcg MusicsFcg
		if err := json.Unmarshal(data, &musicFcg); err == nil {
			if musicFcg.Code == 0 && musicFcg.Req0.Code == 0 {
				if len(musicFcg.Req0.Data.Sip) > 0 && len(musicFcg.Req0.Data.Midurlinfo) > 0 {
					return musicFcg.Req0.Data.Sip[0] + musicFcg.Req0.Data.Midurlinfo[0].PUrl
				}
			}
		}
	}
	return ""
	//return "https://ws.stream.qqmusic.qq.com/C400003mAan70zUy5O.m4a?guid=6509348615&vkey=7087C69919CC3CBD66B44E1BD541A0AB765E23DF30CFD7D1D72D4477F291049A9EB78D553B38640D175FFE9C1BB2F4522BBBDED6DAF39AB0&uin=0&fromtag=3&r=8165635736102514"
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
