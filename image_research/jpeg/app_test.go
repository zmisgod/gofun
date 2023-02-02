package jpeg

import (
	"fmt"
	"log"
	"testing"
)

func TestNewFile(t *testing.T) {
	obj, err := NewFile("./test.jpeg")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(obj.exportHtml())
}

func TestUtil(t *testing.T) {
	//_chunk := make([]byte, 0)
	//_chunk = append(_chunk, 0, 4, 5, 6, 3, 7, 2, 1)
	//var num int
	//for i := 0; i < 6; i++ {
	//	num, _chunk = shift(_chunk)
	//	fmt.Println(num, _chunk)
	//}
	//fmt.Println(_chunk)
	//str := "00000001"
	//fmt.Println(exCount(str))
	//fmt.Println(numberToString(1123213))
}
