package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var url = "http://www.alfuli.com/fuliba"
var responseBody string

func main() {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	responseBody = string(body)
	var re = regexp.MustCompile(`<article .*>(.*?)<\/article>`)
	if len(re.FindString(responseBody)) > 0 {
		var responseBody = re.FindString(responseBody)
		var rest = regexp.MustCompile(`[<]a\starget=\"_blank\"\shref=\"([a-zA-z]+://[^\s]*)\"\stitle=\"(.*?)\"[>]{1}(.*?)</a>`)
		var result = rest.FindAllString(responseBody, -1)
		for _, v := range result {
			//url
			var urlRegexp = regexp.MustCompile(`([a-zA-z]+:[^\s]*)`)
			var url = strings.Replace(urlRegexp.FindString(v), "\"", "", -1)
			//file name
			fileNameRegexp, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			src := strings.Replace(fileNameRegexp.ReplaceAllString(v, "\n"), " ", "", -1)
			println(src)
			println(url)
		}
	}
}
