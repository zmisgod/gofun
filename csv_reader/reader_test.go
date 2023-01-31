package csv_reader

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

type MoveInfo struct {
	Name       string `json:"name"`
	Player     string `json:"player"`
	Open       string `json:"open"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	SiteNumber string `json:"site"`
	Remark     string `json:"remark"`
}

func TestNewCsvReader(t *testing.T) {
	res, err := NewCsvReader("./imax.csv", false)
	if err != nil {
		log.Fatalln(err)
	}
	res.SetCharset(CharsetUTF8)
	resData, err := res.GetValues()
	if err != nil {
		log.Fatalln(err)
	}
	list := make([]MoveInfo, 0)
	for k, v := range resData {
		if k > 0 && len(v) >= 7 {
			list = append(list, MoveInfo{
				Name:       v[0],
				Player:     v[1],
				Open:       v[2],
				Width:      v[3],
				Height:     v[4],
				SiteNumber: v[5],
				Remark:     v[6],
			})
		}
	}
	fmt.Printf("%+v", list)
}

func TestExportUniqueData(t *testing.T) {
	inputName := "./csv_input"
	outputName := "./csv_output"
	list, err := os.ReadDir(inputName)
	if err != nil {
		t.Fatal(err)
	}
	rowMap := make(map[string]string)
	for _, v := range list {
		res, err := NewCsvReader(inputName+"/"+v.Name(), false)
		if err != nil {
			log.Fatalln(err)
		}
		res.SetCharset(CharsetUTF8)
		resData, err := res.GetValues()
		if err != nil {
			t.Fatal(err)
		}
		for _, j := range resData {
			rowMap[j[0]] = j[0]
		}
	}
	//if err := exportCsv(outputName, rowMap); err != nil {
	//	t.Fatal(err)
	//}
	if err := exportQuote(outputName, rowMap); err != nil {
		t.Fatal(err)
	}
}

func exportCsv(outputName string, rowMap map[string]string) error {
	_file, err := os.Create(outputName + "/" + fmt.Sprintf("%s.csv", time.Now().Format("2006-01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = _file.Close()
	}()
	_w := csv.NewWriter(_file)
	for _, v := range rowMap {
		oneRecord := make([]string, 0)
		oneRecord = append(oneRecord, v)
		if err := _w.Write(oneRecord); err != nil {
			return err
		}
	}
	_w.Flush()
	return nil
}

func exportQuote(outputName string, rowMap map[string]string) error {
	_file, err := os.Create(outputName + "/" + fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = _file.Close()
	}()
	rows := make([]string, 0)
	for _, v:=range rowMap {
		rows = append(rows, v)
	}
	str := strings.Join(rows, ",")
	_, err = _file.WriteString(str)
	return err
}
