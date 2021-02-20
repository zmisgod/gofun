package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/axgle/mahonia"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB
var ampSec string
var fileName = "line_lists.csv"

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

func main() {
	trainLists := make(map[string]string)
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := mahonia.NewDecoder("gbk")
	r := csv.NewReader(decoder.NewReader(file))
	reads, err := r.ReadAll()
	checkError(err)

	rowCount := len(reads)
	for i := 0; i < rowCount; i++ {
		trainLists[reads[i][0]] = reads[i][0]
	}
	for _, v := range trainLists {
		searchLineExists(v)
	}
}

func searchLineExists(lineName string) bool {
	var counts int
	err := dbCon.QueryRow(fmt.Sprintf("select count(train_id) as counts from crh_line_lists where train_name = '%s'", lineName)).Scan(&counts)
	if err != nil {
		return false
	}
	if counts > 0 {
		return false
	}
	stmt, err := dbCon.Prepare("insert into crh_line_lists (train_name, type, updated_at, disabled) values (?,?,?,?)")
	if err != nil {
		return false
	}
	_, err = stmt.Exec(lineName, 1, time.Now().Unix(), 0)
	if err != nil {
		return false
	}
	return true
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
