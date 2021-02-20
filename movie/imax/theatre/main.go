package main
//获取各省I-MAX电影院
import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var targetURL = "http://www.imax.cn/theatres-gettheatres"

type result struct {
	ID          string `json:"id"`
	NameEN      string `json:"theatre_name_en"`
	NameCN      string `json:"theatre_name_chinese"`
	SourceID    string `json:"cinemasource_id"`
	Country     string `json:"country_chinese"`
	AddressCN   string `json:"address_chinese_1"`
	AddressEN   string `json:"address_en"`
	Phone       string `json:"contact_phone"`
	URL         string `json:"url_website"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	TheatreType string `json:"theatre_type"`
	VideoType   string `json:"video_type"`
	State       string `json:"state"`
	Province    string `json:"province_name"`
	City        string `json:"city_name"`
}

var dbCon *sql.DB

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func init(){
	err := godotenv.Load("./../../.env")
	checkError(err)
	dbHost := os.Getenv("mysql.host")
	dbPort := os.Getenv("mysql.port")
	dbUser := os.Getenv("mysql.username")
	dbPass := os.Getenv("mysql.password")
	dbName := os.Getenv("mysql.dbname")
	dbCon, err = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	checkError(err)
	defer dbCon.Close()
}

func main() {
	provinceList := []string{
		"北京市",
		"上海市",
		"天津市",
		"重庆市",
		"黑龙江省",
		"吉林省",
		"辽宁省",
		"内蒙古",
		"河北省",
		"新疆",
		"甘肃省",
		"青海省",
		"陕西省",
		"宁夏",
		"河南省",
		"山东省",
		"山西省",
		"安徽省",
		"湖北省",
		"湖南省",
		"江苏省",
		"四川省",
		"贵州省",
		"云南省",
		"广西省",
		"西藏",
		"浙江省",
		"江西省",
		"广东省",
		"福建省",
		"台湾省",
		"海南省",
		"香港",
		"澳门",
	}
	for _,province :=range provinceList {
		resp, err := http.PostForm(targetURL, url.Values{
			"nowlat": {"33.013797"},
			"nowlng": {"119.368489"},
			"city":   {province},
		})
		defer resp.Body.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			rd := bufio.NewReader(resp.Body)
			byteBody, err := ioutil.ReadAll(rd)
			if err != nil {
				fmt.Println(err)
			} else {
				data := bytes.TrimPrefix(byteBody, []byte("\xef\xbb\xbf"))
				var res []result
				err := json.Unmarshal(data, &res)
				if err != nil {
					fmt.Println(err)
				} else {
					for i := 0; i < len(res); i++ {
						dbCon.Exec()
					}
				}
			}
		}
	}
}
