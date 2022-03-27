package novel

import (
	"context"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

func getTextFile(textName string) string {
	return fmt.Sprintf("%s.txt", textName)
}

func FetchChapter(ctx context.Context, textName, url string) error {
	header := map[string]string{
		"sec-fetch-user": "?1",
		"sec-fetch-dest": "document",
		"sec-fetch-mode": "navigate",
		"sec-fetch-site": "none",
		"user-agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	}
	resp, err := utils.HttpClient(ctx, utils.HttpClientMethodGet, url, time.Duration(5)*time.Second, "", nil, header)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body := resp.Body
	document, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return err
	}
	content := document.Find("#newlist ul li")
	content.EachWithBreak(func(i int, contentSelection *goquery.Selection) bool {
		link, _ := contentSelection.Find("a").Attr("href")
		chapterName, _ := contentSelection.Find("a").Html()
		var re = regexp.MustCompile(`(?m)第([0-9]*)章`)
		number := re.FindAllStringSubmatch(chapterName, 10000)
		nameL := strings.Split(chapterName, " ")
		if len(number) > 0 && len(number[0]) >= 2 && len(nameL) > 1 {
			realName := strings.Builder{}
			for k, v := range nameL {
				if k != 0 {
					realName.WriteString(v)
				}
			}
			fmt.Println(link, number[0][1], realName.String())
		}
		//chatp := "第" + strconv.Itoa(j) + "章"
		//fileName := textName + "/" + chatp + ".txt"
		//utils.CreateFile(fileName)
		//file, err := os.OpenFile(getTextFile(textName), os.O_APPEND|os.O_WRONLY, 0777)
		//utils.CheckError(err)
		//if strings.Contains(chapterName, chatp) {
		//	go fetchFromServer(url+"/"+link, chapterName, file)
		//}
		return true
	})
	return nil
}

func fetchFromServer(url, chapterName string, fp *os.File) {
	document, err := goquery.NewDocument(url)
	utils.CheckError(err)
	content, err := document.Find("#content").Html()
	utils.CheckError(err)
	decoder := mahonia.NewDecoder("gbk")
	content = strings.Replace(content, "&nbsp;", "", len(content))
	content = strings.Replace(content, "<br/>", "", len(content))
	content = decoder.ConvertString(content)
	content = strings.Replace(content, "聽", "", len(content))
	if content == "" {
		unavailableStr, _ := document.Find("title").Html()
		//如果遇到服务不可用503，直接等待10秒后重试
		if strings.Contains(unavailableStr, "503") {
			time.Sleep(time.Second)
			fetchFromServer(url, chapterName, fp)
		} else {
			log.Println("unexpected error")
			os.Exit(0)
		}
	} else {
		content = "\n\n" + chapterName + "\n\n" + content
		io.WriteString(fp, content)
		log.Printf("%s download finish\n", chapterName)
	}
}
