package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB

//Detail 详细数据
type Detail struct {
	StartStationName string `json:"start_station_name"`
	ArriveTime       string `json:"arrive_time"`
	StationTrainCode string `json:"station_train_code"`
	StationName      string `json:"station_name"`
	TrainClassName   string `json:"train_class_name"`
	ServiceType      string `json:"service_type"`
	StartTime        string `json:"start_time"`
	StopoverTime     string `json:"stopover_time"`
	EndStationName   string `json:"end_station_name"`
	StationNo        string `json:"station_no"`
	IsEnabled        bool   `json:"isEnabled"`
}

//TwoLevel 二级数据
type TwoLevel struct {
	Data []Detail `json:"data"`
}

//OneLevel 一级数据
type OneLevel struct {
	ValidateMessagesShowID string   `json:"validateMessagesShowId"`
	Status                 bool     `json:"status"`
	Httpstatus             int      `json:"httpstatus"`
	Data                   TwoLevel `json:"data"`
}

//DBGroup group数据
type DBGroup struct {
	ID            int    `json:"id"`
	TrainNo       string `json:"train_no"`
	FirstTelecode string `json:"first_telecode"`
	EndTelecode   string `json:"end_telecode"`
}

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
	stop := false
	pageSize := 500
	page := 1
	for stop == false {
		res, err := dbCon.Query(fmt.Sprintf("select id,train_no,first_telecode,end_telecode from crh_train_list where is_update = 0 order by id desc limit %d, %d", (page-1)*pageSize, pageSize))
		checkError(err)
		defer res.Close()
		j := 0
		for res.Next() {
			var group DBGroup
			err = res.Scan(&group.ID, &group.TrainNo, &group.FirstTelecode, &group.EndTelecode)
			checkError(err)
			if group.FirstTelecode != "" && group.EndTelecode != "" && group.TrainNo != "" {
				response, err := http.Get(fmt.Sprintf("https://kyfw.12306.cn/otn/czxx/queryByTrainNo?train_no=%s&from_station_telecode=%s&to_station_telecode=%s&depart_date=2018-08-06", group.TrainNo, group.FirstTelecode, group.EndTelecode))
				checkError(err)
				defer response.Body.Close()
				body, err := ioutil.ReadAll(response.Body)
				checkError(err)
				parseJSON(body, group.ID)
			}
			j++
		}
		page++
		if j == 0 {
			stop = true
		}
	}
}

//解析json并且保存数据库
func parseJSON(res []byte, groupID int) {
	var oneLevel OneLevel
	json.Unmarshal(res, &oneLevel)
	trainClassName := ""
	serviceType := ""
	if len(oneLevel.Data.Data) > 0 {
		for i := 0; i < len(oneLevel.Data.Data); i++ {
			if i == 0 {
				trainClassName = oneLevel.Data.Data[i].TrainClassName
				serviceType = oneLevel.Data.Data[i].ServiceType
				updatelistGroup(trainClassName, serviceType, groupID)
			}
			insertListGroup(oneLevel.Data.Data[i], groupID)
		}
	} else {
		updateIsUpdatedGroup(groupID)
	}
}

//updatelistGroup 修改crh_train_group表
func updatelistGroup(trainClassName, serviceType string, groupID int) {
	stmt, err := dbCon.Prepare("update crh_train_list set service_type =? , train_class_name =?,is_update = ? where id = ?")
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(serviceType, trainClassName, 1, groupID)
	checkError(err)
}

//updateIsUpdatedGroup 修改crh_train_group表
func updateIsUpdatedGroup(groupID int) {
	stmt, err := dbCon.Prepare("update crh_train_list set is_update =? where id = ?")
	defer stmt.Close()
	checkError(err)
	_, err = stmt.Exec(1, groupID)
	checkError(err)
}

//crh_train_group_details 新增crh_train_group表
func insertListGroup(detail Detail, groupID int) {
	stmt, err := dbCon.Prepare("insert into crh_train_list_details (group_id,arrive_time,station_name, start_time, stopover_time,station_no, isEnabled) value (?,?,?,?,?,?,?)")
	defer stmt.Close()
	checkError(err)
	isEnable := 0
	if detail.IsEnabled {
		isEnable = 1
	}
	_, err = stmt.Exec(groupID, detail.ArriveTime, detail.StationName, detail.StartTime, detail.StopoverTime, detail.StationNo, isEnable)
	checkError(err)
}

//错误检查
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
