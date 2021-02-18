package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var url = ""
var limit = 10278
var saveFilePath = "/Users/meow/Documents/svg/"

func main() {
	for i := 10278; i <= limit; i++ {
		var ln = strconv.Itoa(i)
		status := make(chan string)
		go fetchPage(ln, url+ln, status)
		response := <-status
		if response != "" {
			fmt.Println(url + ln + " ---- " + response)
		}
	}
}

func fetchPage(ln, url string, nowStatus chan string) {
	res, err := goquery.NewDocument(url)
	if err != nil {
		nowStatus <- err.Error()
		return
	}
	contents, err := res.Find("script").Html()
	if err != nil {
		nowStatus <- err.Error()
		return
	}
	var rest string
	realHTML := strings.Replace(html.UnescapeString(contents), " ", "", 100000)
	urlRs := parseJSON(realHTML)
	if urlRs != "" {
		status := make(chan string)
		go saveFiles(ln, urlRs, status)
		rest = <-status
	} else {
		rest = "do not save anything"
	}
	nowStatus <- rest
}

func saveFiles(name, url string, status chan string) {
	resp, err := http.Get(url)
	if err != nil {
		status <- err.Error()
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		status <- err.Error()
		return
	}
	file, err := os.Create(saveFilePath + name + ".svg")
	if err != nil {
		status <- err.Error()
		return
	}
	defer file.Close()
	_, err = file.Write(body)
	if err != nil {
		status <- err.Error()
	}
	status <- name + " _ok"
}

//IconNormal S
type IconNormal struct {
	URL string `json:"url"`
}

//LOGOINFOOBJ S
type LOGOINFOOBJ struct {
	IconNormal IconNormal `json:"icon_normal"`
}
type editFont struct {
	Editfont LOGOINFOOBJ `json:"editFontInfo"`
}

//parseJson 解析json
func parseJSON(JsVAR string) string {
	rows := map[string]string{}
	res := strings.Split(html.UnescapeString(JsVAR), ";")
	jsonArr := make([]string, 0)
	for _, v := range res {
		if len(v) != 0 {
			jsonArr = append(jsonArr, v)
		}
	}
	for _, v := range jsonArr {
		res := strings.SplitN(v, "=", 2)
		if len(res) == 2 {
			rows[res[0]] = res[1]
		}
	}
	if len(rows) == 0 {
		return ""
	}
	for _, v := range rows {
		var urlArr editFont
		err := json.Unmarshal([]byte(v), &urlArr)
		if err == nil {
			return urlArr.Editfont.IconNormal.URL
		}
	}
	return ""
}
