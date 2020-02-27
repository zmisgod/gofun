package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/axgle/mahonia"
)

func main() {
	trainLists := make(map[string]string)
	file, err := os.Open("alert.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	decoder := mahonia.NewDecoder("gbk")
	r := csv.NewReader(decoder.NewReader(file))
	reads, err := r.ReadAll()
	checkError(err)
	lastCount := len(reads)

	rowCount := len(reads[0])
	fmt.Println(rowCount)
	for i := 0; i < rowCount; i++ {
		trainLists[reads[0][i]] = reads[0][i]
	}
	fmt.Println(len(trainLists))
	fmt.Println(trainLists)
	// for k, v := range reads {
	// 	fmt.Println(k)
	// 	fmt.Println(v)
	// }
	fmt.Println(lastCount)
}
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
