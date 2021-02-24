package music

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type QQMusic struct {
	musicId    string
	linkInfo   *LinkInfo
	qqMusicObj *qqMusicObj
}

type qqMusicObj struct {
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

func (a *QQMusic) Fetch(ctx context.Context, linkInfo *LinkInfo) (string, error) {
	a.linkInfo = linkInfo
	link := a.getQQMusicUrl(ctx)
	if link == "" {
		return "", ErrorLinkIsEmpty
	}
	return link, nil
}

type musicsFcg struct {
	Code uint `json:"code"`
	Req0 struct {
		Code uint `json:"code"`
		Data struct {
			Midurlinfo []struct {
				PUrl string `json:"purl"`
			} `json:"midurlinfo"`
			Sip []string `json:"sip"`
		} `json:"data"`
	} `json:"req_0"`
}

func (a *QQMusic) getQQMusicUrl(ctx context.Context) string {
	songId := a.getMusicID(a.linkInfo.Url)
	if songId == "" {
		return ""
	}
	_url := fmt.Sprintf("https://y.qq.com/n/yqq/song/%s.html?ADTAG=h5_playsong&no_redirect=1", songId)
	a.getMusicInfo(_url)
	if a.qqMusicObj == nil {
		return ""
	}
	getSign, dataStr, qq := a.getQQMusicRequestSign(ctx)
	rest := "https://u.y.qq.com/cgi-bin/musics.fcg?-=getplaysongvkey3047664478189791&g_tk=743553847&sign=" + getSign + "&loginUin=" + strconv.Itoa(qq) + "&hostUin=0&format=json&inCharset=utf8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=0&data=" + url.PathEscape(dataStr)
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
	var musicFcg musicsFcg
	if err := json.Unmarshal(data, &musicFcg); err == nil {
		if musicFcg.Code == 0 && musicFcg.Req0.Code == 0 {
			if len(musicFcg.Req0.Data.Sip) > 0 && len(musicFcg.Req0.Data.Midurlinfo) > 0 {
				return musicFcg.Req0.Data.Sip[0] + musicFcg.Req0.Data.Midurlinfo[0].PUrl
			}
		}
	}
	return ""
}

func (a *QQMusic) getQQMusicRequestSign(ctx context.Context) (string, string, int) {
	qq := 83881690
	req := fmt.Sprintf("{\"req\":{\"module\":\"CDN.SrfCdnDispatchServer\",\"method\":\"GetCdnDispatch\",\"param\":{\"guid\":\"6509348615\",\"calltype\":0,\"userip\":\"\"}},\"req_0\":{\"module\":\"vkey.GetVkeyServer\",\"method\":\"CgiGetVkey\",\"param\":{\"guid\":\"6509348615\",\"songmid\":[\"%s\"],\"songtype\":[0],\"uin\":\"%d\",\"loginflag\":%d,\"platform\":\"20\"}},\"comm\":{\"uin\":%d,\"format\":\"json\",\"ct\":24,\"cv\":0}}", a.qqMusicObj.SongMId, qq, 1, qq)
	sign := a.FetchQQMusicSign(ctx, req)
	if sign == "" {
		return "", "", 0
	}
	return sign, req, qq
}

func (a *QQMusic) FetchQQMusicSign(ctx context.Context, body string) string {
	postData := make(url.Values)
	postData["needSign"] = []string{body}
	resp, err := http.PostForm(QQMusicGetSignUrl, postData)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	rest, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	type nodeResult struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}
	var node nodeResult
	_ = json.Unmarshal(rest, &node)
	if node.Code == 200 {
		return node.Data
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

func (a *QQMusic) getMusicInfo(urlStr string) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	doc.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		str, _ := selection.Html()
		if str != "" {
			str = parseUrl(str)
			var re = regexp.MustCompile(`var g_SongData = (.*);`)
			list := re.FindAllStringSubmatch(str, 1000)
			if len(list) > 0 {
				if len(list[0]) >= 2 {
					_ = json.Unmarshal([]byte(list[0][1]), &a.qqMusicObj)
				}
			}
		}
		return true
	})
}

func (a *QQMusic) getMusicID(urlStr string) string {
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
	songId := ""
	if resultUrl != "" {
		var re = regexp.MustCompile(`(?m)&mid=(.*)&`)
		list := re.FindAllStringSubmatch(resultUrl, 1000)
		if len(list) > 0 {
			if len(list[0]) >= 2 {
				songId = list[0][1]
			}
		}
	}
	return songId
}
