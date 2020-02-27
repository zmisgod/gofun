package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB

func init() {
	err := godotenv.Load("./../../.env")
	checkError(err)
	dbHost := os.Getenv("mysql.host")
	dbPort := os.Getenv("mysql.port")
	dbUser := os.Getenv("mysql.username")
	dbPass := os.Getenv("mysql.password")
	dbName := os.Getenv("mysql.dbname")
	dbCon, err = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	checkError(err)
}

//StationInfo StationInfo
type StationInfo struct {
	CnName       string  `json:"cn_name"`
	Longtitude   string  `json:"longtitude"`
	Latitude     string  `json:"latitude"`
	Next         string  `json:"next"`
	NextDistance float64 `json:"next_distance"`
}

//CompareObj CompareObj
type CompareObj struct {
	Station string  `json:"station"`
	Compare float64 `json:"compare"`
}

func main() {
	// var compareResult []CompareObj
	// firstName := "北京南"
	endName := "上海虹桥"
	// sql := "SELECT id,cn_name, ( 3959 * acos( cos( radians('39.865208') ) * cos( radians( latitude ) ) * cos( radians( longtitude ) - radians('116.378596') ) + sin( radians('39.865208') ) * sin( radians( latitude ) ) ) ) AS distance FROM crh_stations where id in (4,698,48,1650,528,1639,565,57,14,799,1712,1919,2073,1781,1638,1429,691,18,252,692,1163,793,789) HAVING distance < 10000000 ORDER BY distance LIMIT 0, 500"
	sql := "select cn_name, longtitude, latitude from crh_stations where id in (4,698,48,1650,528,1639,565,57,14,799,1712,1919,2073,1781,1638,1429,691,18,252,692,1163,793,789)"
	res, err := dbCon.Query(sql)
	checkError(err)
	defer res.Close()
	var stationInfoLists []StationInfo
	for res.Next() {
		var stationInfo StationInfo
		err = res.Scan(&stationInfo.CnName, &stationInfo.Longtitude, &stationInfo.Latitude)
		checkError(err)
		stationInfoLists = append(stationInfoLists, stationInfo)
	}
	//第一个与每个车站信息比较，选出最近的一个
	for index, search := range stationInfoLists {
		if search.CnName == endName {
			continue
		}
		var compareObjLists []CompareObj
		for _, two := range stationInfoLists {
			if two.CnName != search.CnName {
				lat1, _ := strconv.ParseFloat(search.Latitude, 64)
				lon1, _ := strconv.ParseFloat(search.Longtitude, 64)

				lat2, _ := strconv.ParseFloat(two.Latitude, 64)
				lon2, _ := strconv.ParseFloat(two.Longtitude, 64)
				re := EarthDistance(lat1, lon1, lat2, lon2)
				var compareObj CompareObj
				compareObj.Compare = re
				compareObj.Station = two.CnName
				compareObjLists = append(compareObjLists, compareObj)
			}
		}
		last := 0.0
		var thisResult CompareObj
		for k, com := range compareObjLists {
			if k == 0 {
				last = com.Compare
				thisResult.Compare = com.Compare
				thisResult.Station = com.Station
			}
			if com.Compare < last {
				last = com.Compare
				thisResult.Compare = com.Compare
				thisResult.Station = com.Station
			}
		}
		stationInfoLists[index].Next = thisResult.Station
		stationInfoLists[index].NextDistance = thisResult.Compare
	}
	fmt.Println(stationInfoLists)
}

//EarthDistance 两点之间距离
func EarthDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6371000.0 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
