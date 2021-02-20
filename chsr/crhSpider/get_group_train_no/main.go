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

//TrainGroup group
type TrainGroup struct {
	FirstStation string `json:"first_station"`
	EndStation   string `json:"end_station"`
	ID           int    `json:"id"`
}

//TrainList list
type TrainList struct {
	TrainNo string `json:"train_no"`
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
	for stop == false {
		sql := fmt.Sprintf("select first_station,end_station,id from crh_train_group where train_no = ''")
		res, err := dbCon.Query(sql)
		defer res.Close()
		count := 0
		checkError(err)
		for res.Next() {
			var trainName TrainGroup
			err := res.Scan(&trainName.FirstStation, &trainName.EndStation, &trainName.ID)
			checkError(err)
			trainNo := getTrainNo(trainName.FirstStation, trainName.EndStation)
			if trainNo != "" {
				updateTrainNo(trainNo, trainName.ID)
			}
			count++
		}
		if count == 0 {
			stop = true
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

//获取train_no
func getTrainNo(firstStation, endStation string) string {
	sql := fmt.Sprintf("select train_no from crh_train_list where first_station = '%s' and end_station ='%s' order by id asc limit 1", firstStation, endStation)
	res, err := dbCon.Query(sql)
	checkError(err)
	defer res.Close()
	for res.Next() {
		var list TrainList
		err := res.Scan(&list.TrainNo)
		checkError(err)
		return list.TrainNo
	}
	return ""
}

func updateTrainNo(trainNo string, id int) {
	stmt, err := dbCon.Prepare("update crh_train_group set train_no = ? where id = ?")
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(trainNo, id)
	checkError(err)
}
