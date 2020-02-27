package main

import (
	"database/sql"
	"fmt"
	"os"

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

//SName 名称
type SName struct {
	StationName string `json:"station_name"`
}

//GroupName group name
type GroupName struct {
	ID           int    `json:"id"`
	FirstStation string `json:"first_station"`
	EndStation   string `json:"end_station"`
}

func main() {
	stop := false
	pageSize := 500
	i := 1
	for stop == false {
		count := 0
		sql := fmt.Sprintf("select id,first_station,end_station from crh_train_group where id > 1511 limit %d, %d", (i-1)*pageSize, pageSize)
		res, err := dbCon.Query(sql)
		defer res.Close()
		checkError(err)
		for res.Next() {
			var groupName GroupName
			err = res.Scan(&groupName.ID, &groupName.FirstStation, &groupName.EndStation)
			checkError(err)
			sql := fmt.Sprintf("select station_name from crh_train_list_details where group_id in (select id from crh_train_list where first_station = '%s' and end_station = '%s' and service_type != '') group by  station_name", groupName.FirstStation, groupName.EndStation)
			rest, err := dbCon.Query(sql)
			defer rest.Close()
			checkError(err)
			for rest.Next() {
				var sName SName
				err := rest.Scan(&sName.StationName)
				checkError(err)
				stmt, err := dbCon.Prepare("insert into crh_train_group_lines (station_name,group_id) value (?,?)")
				defer stmt.Close()
				checkError(err)
				_, err = stmt.Exec(sName.StationName, groupName.ID)
				checkError(err)
				count++
			}
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
