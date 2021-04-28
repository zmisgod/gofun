package csv_reader

import (
	"fmt"
	"testing"
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
	res, err := NewCsvReader("./imax.csv")
	if err != nil {
		fmt.Println(err)
	} else {
		res.SetCharset(CharsetUTF8)
		resData, err := res.GetValues()
		if err != nil {
			panic(err)
		} else {
			list := make([]MoveInfo, 0)
			for k, v := range resData {
				if k > 0 {
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
	}
}
