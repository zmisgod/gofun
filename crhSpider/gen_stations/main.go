package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//Draw 画的数据结构
type Draw struct {
	FillColor          string
	Stroke             int
	Resize             int
	LongRe             int
	LaRe               int
	LineWeight         int
	Radius             int
	Width              int
	Height             int
	Content            string
	JSONFloder         string
	ViewBoxStartWidth  int
	ViewBoxStartHeight int
	OutJSONName        string
	OutGroupJSONName   string
}

//TrainGroupData group
type TrainGroupData struct {
	ID             int
	FirstStation   string
	EndStation     string
	TrainNo        string
	ServiceType    string
	TrainClassName string
	LineData       []LineData
}

//LineData 画线的数据结构
type LineData struct {
	Latitude     string
	Longtitude   string
	ID           int
	GroupID      int
	ArriveTime   string
	StationName  string
	StartTime    string
	StopoverTime string
	StationNo    int
	ServiceType  int
	IsEnabled    int
}

//CircleData 画圆的数据结构
type CircleData struct {
	Latitude   string
	Longtitude string
	CnName     string
}

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

func main() {
	var draw Draw
	draw.FillColor = "red"
	draw.Stroke = 3
	draw.Resize = 300
	draw.LongRe = 70
	draw.LaRe = 18
	draw.LineWeight = 5
	draw.Radius = 4
	draw.ViewBoxStartWidth = 0
	draw.ViewBoxStartHeight = 0
	draw.OutJSONName = "./stations.svg"
	draw.JSONFloder = "./crhJSON/"
	draw.OutGroupJSONName = "./group_list.json"
	//生成json
	// _, err := draw.createFolder()
	// checkError(err)
	// draw.getLineGroup()

	//生成车站点
	circleList := getStations()
	draw.DrawCircle(circleList)
	draw.Ouput()

	// draw.genTrainLists()

	// draw.DrawAll()
}

//DrawAll 画出所有的线路
func (draw *Draw) DrawAll() {
	circleList := getStations()
	draw.DrawCircle(circleList)
	draw.DrawPolyLines()
	draw.Ouput()
}

//DrawPolyLines 画出所有线路
func (draw *Draw) DrawPolyLines() {
	sql := fmt.Sprintf("select id,first_station,end_station,train_no,service_type,train_class_name from crh_train_group where service_type !='' order by id desc")
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	for res.Next() {
		var trainGroup TrainGroupData
		err := res.Scan(&trainGroup.ID, &trainGroup.FirstStation, &trainGroup.EndStation, &trainGroup.TrainNo, &trainGroup.ServiceType, &trainGroup.TrainClassName)
		checkError(err)
		sql := fmt.Sprintf("select d.id,d.group_id,d.arrive_time,d.station_name,d.start_time,d.stopover_time,d.station_no,d.isEnabled,s.longtitude,s.latitude from crh_train_group_details as d join crh_stations as s on d.station_name = s.cn_name where d.group_id = %d and s.latitude !='' order by d.station_no asc", trainGroup.ID)
		res, err := dbCon.Query(sql)
		defer res.Close()
		checkError(err)
		var lineLists []LineData
		for res.Next() {
			var lineData LineData
			err := res.Scan(&lineData.ID, &lineData.Latitude, &lineData.ArriveTime, &lineData.StationName, &lineData.StartTime, &lineData.StopoverTime, &lineData.StationNo, &lineData.IsEnabled, &lineData.Longtitude, &lineData.Latitude)
			checkError(err)
			lineLists = append(lineLists, lineData)
		}
		svgStr := ""
		if len(lineLists) > 0 {
			svgStr = "<polyline points=\""
			i := 0
			points := ""
			for _, line := range lineLists {
				longtitute := draw.getLongValidNumber(line.Longtitude)
				latitude := draw.getLaValidNumber(line.Latitude)
				longStr := strconv.Itoa(longtitute)
				latiStr := strconv.Itoa(latitude)
				points += ", " + longStr + "," + latiStr + " "
				i++
			}
			points = strings.TrimLeft(points, ", ")
			svgStr += points + "\" fill=\"none\" stroke=\"#000\" stroke-width=\"" + strconv.Itoa(draw.LineWeight) + "\" />"
		}
		draw.Content += svgStr
		checkError(err)
	}
}

//TrainGroupInfo group信息
type TrainGroupInfo struct {
	firstStation   string
	endStation     string
	Name           string `json:"name"`
	TrainNo        string `json:"train_no"`
	TrainClassName string `json:"train_class_name"`
}

func (draw *Draw) genTrainLists() {
	sql := fmt.Sprintf("select first_station,end_station,train_no,train_class_name from crh_train_group where service_type !='' order by id desc")
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	var trainGroupInfoLists []TrainGroupInfo
	for res.Next() {
		var trainGroupInfo TrainGroupInfo
		err := res.Scan(&trainGroupInfo.firstStation, &trainGroupInfo.endStation, &trainGroupInfo.TrainNo, &trainGroupInfo.TrainClassName)
		checkError(err)
		trainGroupInfo.Name = trainGroupInfo.firstStation + " ~ " + trainGroupInfo.endStation
		trainGroupInfoLists = append(trainGroupInfoLists, trainGroupInfo)
	}
	data, err := json.Marshal(trainGroupInfoLists)
	checkError(err)
	createFile(draw.OutGroupJSONName, string(data))
}

//getLineGroup 获取车站列表
func (draw *Draw) getLineGroup() {
	sql := fmt.Sprintf("select id,first_station,end_station,train_no,service_type,train_class_name from crh_train_group where service_type !='' order by id desc")
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	for res.Next() {
		var trainGroup TrainGroupData
		err := res.Scan(&trainGroup.ID, &trainGroup.FirstStation, &trainGroup.EndStation, &trainGroup.TrainNo, &trainGroup.ServiceType, &trainGroup.TrainClassName)
		checkError(err)
		draw.getLineGroupDetail(trainGroup)
	}
}

//getStations 获取车站点
func getStations() []CircleData {
	sql := fmt.Sprintf("select longtitude,latitude,cn_name from crh_stations where latitude != '' order by id desc")
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	var circleList []CircleData
	for res.Next() {
		var cicleData CircleData
		err := res.Scan(&cicleData.Longtitude, &cicleData.Latitude, &cicleData.CnName)
		checkError(err)
		circleList = append(circleList, cicleData)
	}
	return circleList
}

//getLineGroupDetail
func (draw *Draw) getLineGroupDetail(groupData TrainGroupData) {
	sql := fmt.Sprintf("select d.id,d.group_id,d.arrive_time,d.station_name,d.start_time,d.stopover_time,d.station_no,d.isEnabled,s.longtitude,s.latitude from crh_train_group_details as d join crh_stations as s on d.station_name = s.cn_name where d.group_id = %d and s.latitude !='' order by d.station_no asc", groupData.ID)
	res, err := dbCon.Query(sql)
	defer res.Close()
	checkError(err)
	var lineLists []LineData
	for res.Next() {
		var lineData LineData
		err := res.Scan(&lineData.ID, &lineData.Latitude, &lineData.ArriveTime, &lineData.StationName, &lineData.StartTime, &lineData.StopoverTime, &lineData.StationNo, &lineData.IsEnabled, &lineData.Longtitude, &lineData.Latitude)
		checkError(err)
		lineLists = append(lineLists, lineData)
	}
	if len(lineLists) > 0 {
		draw.DrawPolyLine(groupData, lineLists)
	}
}

//DrawPolyLine 画polyLine
func (draw *Draw) DrawPolyLine(groupData TrainGroupData, lines []LineData) {
	svgStr := ""
	i := 0
	for _, line := range lines {
		longtitute := draw.getLongValidNumber(line.Longtitude)
		latitude := draw.getLaValidNumber(line.Latitude)
		if latitude > draw.Width {
			draw.Width = latitude
		}
		longStr := strconv.Itoa(longtitute)
		latiStr := strconv.Itoa(latitude)
		if i == 0 {
			svgStr += ", " + longStr + "," + latiStr + ""
		} else {
			svgStr += ", " + longStr + "," + latiStr + ""
		}
		i++
	}
	var output OutputJSONObject
	output.Path = svgStr
	output.TrainNo = groupData.TrainNo
	output.Fill = "transparent"
	output.Stroke = strconv.Itoa(draw.Stroke)
	output.StrokeWidth = strconv.Itoa(draw.LineWeight)
	output.Path = strings.TrimLeft(svgStr, ", ")
	ress, err := json.Marshal(output)
	checkError(err)
	draw.OuputLineJSON(draw.JSONFloder+groupData.TrainNo+".json", string(ress))
}

//OutputJSONObject 输出的JSON
type OutputJSONObject struct {
	Path        string `json:"path"`
	TrainNo     string `json:"train_no"`
	Fill        string `json:"fill"`
	Stroke      string `json:"stroke"`
	StrokeWidth string `json:"stroke-width"`
}

//OuputLineJSON  输出json文件
func (draw *Draw) OuputLineJSON(fileName, content string) {
	createFile(fileName, content)
}

//getRValidNumber 获取数
func (draw *Draw) getLaValidNumber(res string) int {
	f, err := strconv.ParseFloat(res, 32)
	checkError(err)
	return int((f - float64(draw.LaRe)) * float64(draw.Resize))
}

//getLValidNumber
func (draw *Draw) getLongValidNumber(res string) int {
	f, err := strconv.ParseFloat(res, 32)
	checkError(err)
	return int((f - float64(draw.LongRe)) * float64(draw.Resize))
}

//DrawCircle 画circle
func (draw *Draw) DrawCircle(circles []CircleData) {
	i := 1
	circleStr := ""
	for _, circle := range circles {
		la := draw.getLaValidNumber(circle.Latitude)
		long := draw.getLongValidNumber(circle.Longtitude)
		if la > draw.Width {
			draw.Width = la
		}
		if long > draw.Height {
			draw.Height = long
		}
		circleStr += fmt.Sprintf("<circle cx=\"%d\" cy=\"%d\" r=\"%d\" aid=\"%d\" stroke=\"%d\" fill=\"%s\" alt=\"%s\" />", long, la, draw.Radius, i, draw.Stroke, draw.FillColor, circle.CnName)
		i++
	}
	draw.Content += "<g>" + circleStr + "</g>"
}

//Ouput 输出文件
func (draw *Draw) Ouput() {
	draw.Width += 500
	draw.Height += 500
	res := fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"%d,%d,%d,%d\" version=\"1.1\"  width=\"100%%\" height=\"100%%\">%s</svg>", 0, 0, draw.Height, draw.Width, draw.Content)
	createFile(draw.OutJSONName, res)
}

//生成文件
func createFile(fileName, content string) {
	file, err := os.Create(fileName)
	checkError(err)
	_, err = file.WriteString(content)
	checkError(err)
}

/**
 * 创建文件夹
 */
func (draw *Draw) createFolder() (bool, error) {
	checkFloderNotExists, err := draw.checkPathIsNotExists()
	if err != nil {
		return false, err
	}
	if checkFloderNotExists {
		err := os.MkdirAll(draw.JSONFloder, 0777)
		if err != nil {
			return false, err
		}
		fmt.Printf("create floder %s successful\n", draw.JSONFloder)
		return true, nil
	}
	return false, err
}

/**
 * 检查文件是否存在
 * 返回true 不存在， false 存在
 */
func (draw *Draw) checkPathIsNotExists() (bool, error) {
	_, err := os.Stat(draw.JSONFloder)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//checkError 检查错误
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
