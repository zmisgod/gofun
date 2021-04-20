package csv_reader

import (
	"fmt"
	"testing"
	"time"
)

type MoveInfo struct {
	Name       string    `json:"name"`
	Player     string    `json:"player"`
	Open       string    `json:"open"`
	Width      string    `json:"width"`
	Height     string    `json:"height"`
	SiteNumber string    `json:"site"`
	Remark     string    `json:"remark"`
	Info       []OneInfo `json:"info"`
}

type OneInfo struct {
	Open   time.Time `json:"open"`
	Width  float32   `json:"width"`
	Height float32   `json:"height"`
	Area   float32   `json:"area"`
	Site   int       `json:"site"`
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
						SiteNumber: v[6],
						Remark:     v[7],
						Info:       nil,
					})
				}
			}
			fmt.Printf("%+v", list)
			//nameFile, err := os.Create("name.log")
			//defer nameFile.Close()
			//if err != nil {
			//	panic(err)
			//}else{
			//	write := bufio.NewWriter(nameFile)
			//	for _,v := range resultArr {
			//		fmt.Fprintln(write, v)
			//	}
			//	_ = write.Flush()
			//}
		}
	}
}