package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB
var fileName = "alertA.csv"

func init() {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}
	dbUser := os.Getenv("mysql.username")
	dbPass := os.Getenv("mysql.password")
	dbName := os.Getenv("mysql.dbname")
	dbHost := os.Getenv("mysql.host")
	dbPort := os.Getenv("mysql.port")
	dbCon, err = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func main() {
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()
	decoder := mahonia.NewDecoder("gbk")
	r := csv.NewReader(decoder.NewReader(file))
	reads, err := r.ReadAll()
	checkError(err)

	//第一行的数据 1,1 , 2,1 , 3,1
	columnCount := len(reads[0])
	trainLists := make(map[string]([]string))
	for i := 0; i < columnCount; i++ {
		trainName := reads[0][i]
		stop := false
		j := 1
		stationLists := make([]string, 0)
		for !stop {
			if len(reads) >= j+1 {
				stationName := reads[j][i]
				if stationName != "" {
					stationLists = append(stationLists, stationName)
					j++
				} else {
					stop = true
				}
			} else {
				stop = true
			}
		}
		trainLists[trainName] = stationLists
	}
	for trainName, lists := range trainLists {
		//检查crh_line_lists的数据是否存在
		trainID := searchLinenameExists(trainName)
		if len(lists) > 0 {
			//将之前的数据删除
			needInsert := clearOld(trainID, len(lists))
			if needInsert {
				for k, stationName := range lists {
					if strings.Trim(stationName, " ") != "@@" {
						tType, sort := checkNowStationSort(lists, k, stationName)
						stationName = strings.Replace(stationName, " ", "", len(stationName))
						stationName = strings.Replace(stationName, "站", "", len(stationName))
						insertData(trainID, stationName, tType, sort)
					}
				}
			} else {
				fmt.Println(trainName + " do not need insert")
			}
		}
	}
}

//检查当前站名在列表中的排名
func checkNowStationSort(lists []string, key int, stationName string) (int, int) {
	atPosition := []int{}
	for k, v := range lists {
		if strings.Trim(v, " ") == "@@" {
			atPosition = append(atPosition, k)
		}
	}
	if len(atPosition) == 0 {
		return 0, key
	}
	searchClosetAt := 1000
	searchKey := 1000
	searchIndex := 0
	for index, position := range atPosition {
		if key-position < searchClosetAt && key-position > 0 {
			searchKey = position
			searchClosetAt = key - position
			searchIndex = index
		}
	}
	if key-searchKey > 0 {
		return searchIndex + 1, key - searchKey
	}
	return 0, key
}

func searchLinenameExists(lineName string) int64 {
	var trainID int64
	err := dbCon.QueryRow(fmt.Sprintf("select train_id from crh_line_lists where train_name = '%s'", lineName)).Scan(&trainID)
	if err != nil {
		stmt, _ := dbCon.Prepare("insert into crh_line_lists (train_name, type, updated_at, disabled) values (?,?,?,?)")
		defer stmt.Close()
		rest, _ := stmt.Exec(lineName, 1, time.Now().Unix(), 0)
		lastInsertID, _ := rest.LastInsertId()
		return lastInsertID
	}
	return trainID
}

//插入线路信息
func insertData(trainID int64, stationName string, tType, sort int) {
	stationID := getStationID(stationName)
	stmt, err := dbCon.Prepare("insert into crh_line_details (train_id, station_id, type, sort, disabled) values (?,?,?,?,?)")
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(trainID, stationID, tType, sort, 0)
	checkError(err)
}

//获取站名id
func getStationID(stationName string) int64 {
	var id int64
	err := dbCon.QueryRow(fmt.Sprintf("select id from crh_line_stations where station_name = '%s'", stationName)).Scan(&id)
	if err != nil {
		stmt, err := dbCon.Prepare("insert into crh_line_stations (station_name,station_address,longtitude,latitude,citycode,adcode,pcode, status,type) values (?,?,?,?,?,?,?,?,?)")
		defer stmt.Close()
		checkError(err)
		rest, err := stmt.Exec(stationName, "", "", "", "", "", "", 0, 0)
		checkError(err)
		stationID, _ := rest.LastInsertId()
		return stationID
	}
	return id
}

//将之前的数据删除
func clearOld(trainID int64, count int) bool {
	//先检查总数，如果总数对，就不删除，跳过，如果总数不对就
	var counts int
	err := dbCon.QueryRow(fmt.Sprintf("select count(1) as counts from crh_line_details where train_id = %d", trainID)).Scan(&counts)
	checkError(err)
	if counts == count {
		return false
	} else {
		stmt, err := dbCon.Prepare("delete from crh_line_details where train_id = ?")
		defer stmt.Close()
		checkError(err)
		_, err = stmt.Exec(trainID)
		checkError(err)
	}
	return true
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
