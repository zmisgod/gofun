package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//Geo 返回的数据
type Geo struct {
	Status string    `json:"status"`
	Pois   []Details `json:"pois"`
}

//Details 详情
type Details struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Location string `json:"location"`
	Pcode    string `json:"pcode"`
	Citycode string `json:"citycode"`
	Adcode   string `json:"adcode"`
	Pname    string `json:"pname"`
	Cityname string `json:"cityname"`
	Adname   string `json:"adname"`
}

var dbCon *sql.DB
var ampSec string

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
	ampSec = os.Getenv("amap.secret")
	dbCon, sqlErr = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	if sqlErr != nil {
		fmt.Println(sqlErr)
		os.Exit(0)
	}
}

//Station 站点名称
type Station struct {
	StationName string `json:"station_name"`
	ID          int    `json:"id"`
}

func main() {
	pageSize := 500
	for i := 2700; i <= 5330; i += pageSize {
		res, err := dbCon.Query(fmt.Sprintf("select station_name,id from crh_line_stations where id >= %d and id < %d and station_address = ''", i, i+pageSize))
		checkError(err)
		defer res.Close()
		for res.Next() {
			var station Station
			err = res.Scan(&station.StationName, &station.ID)
			checkError(err)
			fetchAmap(station.StationName, station.ID)
		}
	}
}

//获取高德地图的信息
func fetchAmap(stationName string, ID int) {
	amapURL := fmt.Sprintf("http://restapi.amap.com/v3/place/text?key=%s&keywords=%s站&types=火车站&city=&children=1&offset=20&page=1&extensions=all", ampSec, stationName)
	response, err := http.Get(amapURL)
	checkError(err)
	defer response.Body.Close()
	reader := bufio.NewReader(response.Body)
	bytebody, err := ioutil.ReadAll(reader)
	checkError(err)
	var amapGeo Geo
	json.Unmarshal(bytebody, &amapGeo)
	if len(amapGeo.Pois) > 0 {
		oneInfo := amapGeo.Pois[0]
		if checkType(oneInfo.Type, stationName, oneInfo.Name) {
			var saveinfo SaveInfo
			long, la := getLocation(oneInfo.Location)
			saveinfo.Longtitude = long
			saveinfo.Latitude = la
			saveinfo.Address = getFullAddress(oneInfo)
			saveinfo.Adcode = oneInfo.Adcode
			saveinfo.Citycode = oneInfo.Citycode
			saveinfo.Pcode = oneInfo.Pcode
			updateTelecode(saveinfo, ID)
		}
	}
}
func getFullAddress(detail Details) string {
	return detail.Pname + detail.Cityname + detail.Adname + detail.Address
}

//LA 纬度
//LO 经度
func getLocation(location string) (string, string) {
	locations := strings.Split(location, ",")
	long := locations[0]
	la := locations[1]
	return long, la
}

//SaveInfo 保存的数据
type SaveInfo struct {
	Address    string
	Longtitude string
	Latitude   string
	Citycode   string
	Adcode     string
	Pcode      string
}

func checkType(dType, stationName, apiName string) bool {
	if dType == "" {
		return false
	}
	dTypes := strings.Split(dType, ";")
	if len(dType) == 0 {
		return false
	}
	isOk := false
	for _, v := range dTypes {
		if v == "火车站" {
			isOk = true
		}
	}
	if isOk {
		if strings.Contains(apiName, stationName) {
			return true
		}
	}
	return false
}

//错误检查
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func updateTelecode(saveInfo SaveInfo, ID int) {
	updateSQL := fmt.Sprintf("update crh_line_stations set station_address = ?, longtitude = ?, latitude = ?, citycode = ?, adcode =?, pcode =? where id = ?")
	stmt, err := dbCon.Prepare(updateSQL)
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(saveInfo.Address, saveInfo.Longtitude, saveInfo.Latitude, saveInfo.Citycode, saveInfo.Adcode, saveInfo.Pcode, ID)
	checkError(err)
}
