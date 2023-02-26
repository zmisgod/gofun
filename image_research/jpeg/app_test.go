package jpeg

import (
	"bytes"
	"fmt"
	"github.com/zmisgod/gofun/image_research"
	"log"
	"testing"
)

func TestNewFile(t *testing.T) {
	obj, err := NewFile("./WechatIMG37.jpeg")
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
	//fmt.Println(strconv.ParseInt("0100", 2, 64))
	//fmt.Println(ycrcb2rgb([]float64{ 244.56513902640177, 124.88035222703407, 131.11964777296595 }))

	input := bytes.NewBufferString("01000000").Bytes()
	fmt.Println(input)
	cursor := 1
	length := 10
	subBuffer := image_research.SliceArr(input, cursor)
	fmt.Println("subBuffer", string(subBuffer))
	tempSubBuffer := make([]byte, length-len(subBuffer))
	tempSubBuffer = image_research.FillBytes(tempSubBuffer, '1')
	fmt.Println("temp", string(tempSubBuffer))
	subBuffer = image_research.Concat(subBuffer, tempSubBuffer, length)
	fmt.Println(string(subBuffer))
}
