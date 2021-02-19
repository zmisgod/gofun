package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/zmisgod/gofun/drawsvg"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var dbCon *sql.DB

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

//List 线路
type List struct {
	TrainID   int
	TrainName string
	Type      int
	MaxGroup  int
	Details   []Detail
}

//Detail 详情
type Detail struct {
	StationID int
	TrainID   int
	Type      int
	Sort      int
	Station   Station
}

//Station 车站信息
type Station struct {
	ID             int
	StationName    string
	StationAddress string
	Longtitude     string
	Latitude       string
	Type           int
}

func main() {
	svg := drawsvg.Create()
	rows, err := dbCon.Query(fmt.Sprintf("select train_id, train_name, type from crh_line_lists"))
	defer rows.Close()
	checkError(err)
	var trainLists []List
	for rows.Next() {
		var list List
		err := rows.Scan(&list.TrainID, &list.TrainName, &list.Type)
		checkError(err)
		trainLists = append(trainLists, list)
	}
	if len(trainLists) == 0 {
		fmt.Println("crh_line_lists no data here")
		os.Exit(0)
	}
	stationLists := make(map[int]Station)
	rows, err = dbCon.Query("select id,station_name, station_address, longtitude, latitude, type from crh_line_stations where station_address != ''")
	checkError(err)
	for rows.Next() {
		var station Station
		rows.Scan(&station.ID, &station.StationName, &station.StationAddress, &station.Longtitude, &station.Latitude, &station.Type)
		checkError(err)
		stationLists[station.ID] = station
	}

	svg.SetCircle(4, 6, "#42526e", "transparent")
	for _, one := range trainLists {
		rows, err := dbCon.Query(fmt.Sprintf("select station_id, train_id,type, sort from crh_line_details where train_id = %d", one.TrainID))
		checkError(err)
		defer rows.Close()

		var dpath drawsvg.Path
		dpath.Aid = strconv.Itoa(one.TrainID)
		dpath.Alt = one.TrainName
		dpath.Fill = "transparent"
		if one.Type == 1 {
			dpath.Stroke = "#CC99FF"
			dpath.StrokeWidth = 60
		} else if one.Type == 2 {
			dpath.Stroke = "#99CCFF"
			dpath.StrokeWidth = 50
		} else {
			dpath.Stroke = "#999999"
			dpath.StrokeWidth = 2
		}

		dpath.FillOpacity = "0.4"

		maxGroup := 0
		var details []Detail
		for rows.Next() {
			var detail Detail
			err = rows.Scan(&detail.StationID, &detail.TrainID, &detail.Type, &detail.Sort)
			checkError(err)
			if maxGroup < detail.Type {
				maxGroup = detail.Type
			}
			details = append(details, detail)
		}
		dpath.MaxGroup = maxGroup + 1
		for i := 0; i < dpath.MaxGroup; i++ {
			var ipaths []drawsvg.IPath
			for _, value := range details {
				if value.Type == i {
					station, ok := stationLists[value.StationID]
					if ok {
						var ipath drawsvg.IPath
						ipath.ID = station.ID
						ipath.Group = station.Type
						ipath.Long = station.Latitude
						ipath.Lat = station.Longtitude
						ipath.Directive = 0
						ipaths = append(ipaths, ipath)
					}
				}
			}
			dpath.PathInfo = ipaths
			svg.SetPath(dpath)
		}
		// trainLists[index].MaxGroup = maxGroup
		// trainLists[index].Details = details
	}
	content := svg.Draw()
	createFile(content)
}

func createFile(content string) {
	file, err := os.Create("stations.svg")
	checkError(err)
	defer file.Close()
	io.WriteString(file, content)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
