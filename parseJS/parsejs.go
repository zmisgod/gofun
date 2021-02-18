package parseJS

import (
	"bufio"
	"fmt"
	"os"
)

func ParseJS(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	//for{
		line, err := rd.ReadSlice(30)
		if err != nil {
			//break
			fmt.Println(err)
		}else{
			fmt.Println(string(line))
		}
	//}
}
