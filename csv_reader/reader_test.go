package csv_reader

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestNewCsvReader(t *testing.T) {
	res, err := NewCsvReader("./nickname.csv")
	if err != nil {
		fmt.Println(err)
	}else{
		res.SetCharset(CharsetUTF8)
		resData, err := res.GetValues()
		if err != nil {
			panic(err)
		}else{
			resultArr := make([]string, 0)
			for _, v := range resData {
				if len(v) > 0 {
					for _, j := range v {
						if j != "" {
							resultArr = append(resultArr, j)
						}
					}
				}
			}
			nameFile, err := os.Create("name.log")
			defer nameFile.Close()
			if err != nil {
				panic(err)
			}else{
				write := bufio.NewWriter(nameFile)
				for _,v := range resultArr {
					fmt.Fprintln(write, v)
				}
				_ = write.Flush()
			}
		}
	}
}