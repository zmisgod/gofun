package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//TrainData 列车名称的数据结构(最终保存在数据库的结构)
type TrainData struct {
	TrainName    string `json:"train_name"`    //拼音缩写
	TrainCode    string `json:"train_code"`    //中文
	FirstStation string `json:"first_station"` //英文缩写
	EndStation   string `json:"end_station"`   //中文拼音
	TrainNo      string `json:"train_no"`      //拼音缩写2
	Status       int    `json:"status"`        //更新状态
	ID           int    `json:"id"`            //主键
}

//TrainList 获取json中的数据结构
type TrainList struct {
	StationTrainCode string `json:"station_train_code"`
	TrainNo          string `json:"train_no"`
}

var dbCon *sql.DB

const jsPrefix = "var"

func init() {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}
	var sqlErr error
	dbUser := os.Getenv("mysql.username")
	dbPass := os.Getenv("mysql.password")
	dbName := os.Getenv("mysql.dbname")
	dbHost := os.Getenv("mysql.host")
	dbPort := os.Getenv("mysql.port")
	dbCon, sqlErr = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	if sqlErr != nil {
		fmt.Println(sqlErr)
		os.Exit(0)
	}
}

func main() {
	jsFile, err := os.Open("./train_list.js")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer jsFile.Close()
	fileBuf := bufio.NewReader(jsFile)
	fileInfos, err := ioutil.ReadAll(fileBuf)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	jsString := string(fileInfos)
	res := parseJSON(jsString)
	if err, ok := res["train_list"]; ok {
		parseJSONResult(res["train_list"])
	} else {
		fmt.Println(err)
	}
}
func parseJSONDetail(jsonResult interface{}) {
	r, ok := jsonResult.(string)
	if ok {
		fmt.Println("HYUGHUY")
		fmt.Println(r)
	} else {
		fmt.Println(ok)
		fmt.Println(r)
	}
}

//parseJSONResult 将json结果转成结构体返回数组的结构体
func parseJSONResult(jsonResult string) {
	var jsonUn map[string](map[string]interface{})
	err := json.Unmarshal([]byte(jsonResult), &jsonUn)
	if err == nil {
		for _, v := range jsonUn {
			for _, trainLists := range v {
				ktArr, ok := trainLists.([]interface{})
				if ok {
					for _, vt := range ktArr {
						obj, ok := vt.(map[string]interface{})
						if ok {
							stationTrainCode := obj["station_train_code"].(string)
							stationTrainNo := obj["train_no"].(string)
							parseTrainDetail(stationTrainNo, stationTrainCode)
						}
					}
				}
			}
		}
	}
}

func parseTrainDetail(trainNo, trainInfo string) {
	var re = regexp.MustCompile(`\((.*?)\)`)

	matches := ""
	for _, match := range re.FindAllString(trainInfo, -1) {
		matches = match
	}
	prefix := strings.Replace(trainInfo, matches, "", 100)
	matches = strings.Replace(matches, "(", "", 100)
	matches = strings.Replace(matches, ")", "", 100)
	stationArr := strings.Split(matches, "-")
	first := stationArr[0]
	end := stationArr[1]
	saveDB(trainNo, prefix, first, end)
}

//parseJson 解析json
func parseJSON(JSVAR string) map[string]string {
	rows := make(map[string]string)
	res := strings.Split(html.UnescapeString(JSVAR), ";")
	jsonArr := make([]string, 0)
	for _, v := range res {
		if len(v) != 0 {
			jsonArr = append(jsonArr, v)
		}
	}
	for _, v := range jsonArr {
		res := strings.SplitN(v, "=", 2)
		if len(res) == 2 {
			date := parsePrefix(res[0])
			rows[date] = res[1]
		}
	}
	if len(rows) == 0 {
		return rows
	}
	return rows
}

//parsePrefix  切分前缀
func parsePrefix(prefix string) string {
	splits := strings.Split(jsPrefix, ",")
	for _, v := range splits {
		res := strings.Split(prefix, v)
		if len(res) == 2 {
			return strings.Replace(res[1], " ", "", 1000)
		}
	}
	return prefix
}

//saveDB 保存到DB
func saveDB(trainNo, trainCode, firstStation, endStation string) {
	res, err := dbCon.Query(fmt.Sprintf("select id from crh_train_list where train_no = '%s' and train_code = '%s'", trainNo, trainCode))
	defer res.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		id := 0
		for res.Next() {
			var ret TrainData
			err = res.Scan(&ret.ID)
			if err != nil {
				fmt.Println(err)
			} else {
				id = ret.ID
			}
		}
		if id == 0 {
			stmt, err := dbCon.Prepare("insert into crh_train_list (train_no, train_code, first_station, end_station,  status) values (?,?,?,?,?)")
			defer stmt.Close()
			if err != nil {
				fmt.Println(err)
			} else {
				_, err = stmt.Exec(trainNo, trainCode, firstStation, endStation, 1)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		// else {
		// 	stmt, err := dbCon.Prepare("update crh_train_list set status = ? where id = ?")
		//  defer stmt.Close()
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	} else {
		// 		_, err = stmt.Exec(1, id)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 		}
		// 	}
		// }
	}
}
