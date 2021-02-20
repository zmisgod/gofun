package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

var dbCon *sql.DB

//DistinctTrain 列车
type DistinctTrain struct {
	FirstStation string `json:"first_station"`
	EndStation   string `json:"end_station"`
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
	pageSize := 500
	i := 1
	for stop == false {
		count := 0
		sql := fmt.Sprintf("select first_station,end_station from crh_train_list group by first_station,end_station limit %d, %d", (i-1)*pageSize, pageSize)
		res, err := dbCon.Query(sql)
		defer res.Close()
		checkError(err)
		for res.Next() {
			var trainInfo DistinctTrain
			err := res.Scan(&trainInfo.FirstStation, &trainInfo.EndStation)
			checkError(err)
			stmt, err := dbCon.Prepare("insert into crh_train_group (first_station, end_station, train_no, first_telecode, end_telecode, status, service_type, train_class_name) value (?,?,?,?,?,?,?,?)")
			defer stmt.Close()
			checkError(err)
			_, err = stmt.Exec(trainInfo.FirstStation, trainInfo.EndStation, "", "", "", 0, "", "")
			checkError(err)
			count++
		}
		if count == 0 {
			stop = true
		}
		i++
	}
	fmt.Println("stop")
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
