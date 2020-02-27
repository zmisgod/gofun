package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB

//PAGESIZE 分页
const PAGESIZE = 500

//TrainGroup city code
type TrainGroup struct {
	FirstStation string `json:"first_station"`
	EndStation   string `json:"end_station"`
	ID           int    `json:"id"`
}

//Telecode city code
type Telecode struct {
	Telecode string `json:"telecode"`
}

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

func main() {
	stop := false
	i := 1
	mapsArr := getTrainNamesLists()
	for stop == false {
		sql := fmt.Sprintf("select first_station,end_station,id from crh_train_group where (end_telecode = '' or first_telecode = '') limit %d, %d", (i-1)*PAGESIZE, PAGESIZE)
		res, err := dbCon.Query(sql)
		defer res.Close()
		count := 0
		checkError(err)
		for res.Next() {
			var trainName TrainGroup
			err := res.Scan(&trainName.FirstStation, &trainName.EndStation, &trainName.ID)
			checkError(err)
			updateSQL := fmt.Sprintf("first_telecode = '%s', end_telecode = '%s'", mapsArr[trainName.FirstStation], mapsArr[trainName.EndStation])
			updateTelecode(trainName.ID, updateSQL, "crh_train_group")
			count++
		}
		if count == 0 {
			stop = true
		}
		i++
	}
}

//StationNames s
type StationNames struct {
	CnName   string `json:"cn_name"`
	Telecode string `json:"telecode"`
}

func getTrainNamesLists() map[string]string {
	saveNames := make(map[string]string)
	sql := fmt.Sprintf("select cn_name,telecode from crh_stations")
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	for res.Next() {
		var names StationNames
		err = res.Scan(&names.CnName, &names.Telecode)
		checkError(err)
		saveNames[names.CnName] = names.Telecode
	}
	return saveNames
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

//fetchTelecode 获取telecode
func fetchTelecode(cityName string) string {
	sql := fmt.Sprintf("select telecode from crh_stations where cn_name = '%s'", cityName)
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	telecode := ""
	for res.Next() {
		var trainInfo Telecode
		err := res.Scan(&trainInfo.Telecode)
		checkError(err)
		telecode = trainInfo.Telecode
		break
	}
	return telecode
}

//updateTelecode 修改telecode
func updateTelecode(id int, updateParams, tableName string) {
	updateSQL := fmt.Sprintf("update %s set %s where id = ?", tableName, updateParams)
	stmt, err := dbCon.Prepare(updateSQL)
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(id)
	checkError(err)
}
