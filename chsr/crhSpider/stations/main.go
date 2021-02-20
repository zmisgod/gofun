package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//TrainData 列车名称的数据结构
type TrainData struct {
	PyAbbreviate    string `json:"abb_name"`  //拼音缩写
	Name            string `json:"cn_name"`   //中文
	Telecode        string `json:"telecode"`  //英文缩写
	NamePy          string `json:"py_name"`   //中文拼音
	PyAbbreviateTwo string `json:"abb_name2"` //拼音缩写2
	OrderNumber     int    `json:"no"`        //序号
	Status          int    `json:"status"`    //更新状态
	ID              int    `json:"id"`        //主键
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
	jsFile, err := os.Open("./station_name.js")
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
	if err, ok := res["station_names"]; ok {
		jsonArr := parseJSONResult(res["station_names"])
		saveDB(jsonArr)
	} else {
		fmt.Println(err)
	}
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
			prefix := parsePrefix(res[0])
			rows[prefix] = res[1]
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

//parseJSONResult 将json结果转成结构体返回数组的结构体
func parseJSONResult(json string) []TrainData {
	res := strings.Split(json, "@")
	var trainDataLists []TrainData
	for _, v := range res {
		var trainData TrainData
		detail := strings.Split(v, "|")
		if len(detail) == 6 {
			trainData.PyAbbreviate = detail[0]
			trainData.Name = detail[1]
			trainData.Telecode = detail[2]
			trainData.NamePy = detail[3]
			trainData.PyAbbreviateTwo = detail[4]
			trainData.OrderNumber, _ = strconv.Atoi(detail[5])
			trainDataLists = append(trainDataLists, trainData)
		}
	}
	return trainDataLists
}

//saveDB 保存到DB
func saveDB(lists []TrainData) {
	for _, v := range lists {
		res, err := dbCon.Query(fmt.Sprintf("select id from crh_stations where telecode = '%s'", v.Telecode))
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
				stmt, err := dbCon.Prepare("insert into crh_stations (abb_name, cn_name, telecode, py_name, abb_name2,no, status) values (?,?,?,?,?,?,?)")
				defer stmt.Close()
				if err != nil {
					fmt.Println(err)
				} else {
					_, err = stmt.Exec(v.PyAbbreviate, v.Name, v.Telecode, v.NamePy, v.PyAbbreviateTwo, v.OrderNumber, 1)
					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				stmt, err := dbCon.Prepare("update crh_stations set status = ? where id = ?")
				defer stmt.Close()
				if err != nil {
					fmt.Println(err)
				} else {
					_, err = stmt.Exec(1, id)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
