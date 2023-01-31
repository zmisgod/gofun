package jpeg

import (
	"fmt"
	"log"
	"testing"
)

func TestNewFile(t *testing.T) {
	data, err := NewFile("./test.jpeg")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(data)
}

func TestUtil(t *testing.T) {
	//1
	//bu := bufferToString([]byte{
	//	0xff, 0xd8, 0xff,0xe0, 0x00,0x10, 0x4a, 0x46,
	//	0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01,
	//})
	//fmt.Println(strToBuffer(bu))
	//res := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	//fmt.Println(sliceArr(res, 6))
	//2
	//_i, _ := strconv.ParseInt("1001", 2, 10)
	//_plus := _i+1
	//fmt.Println(strconv.FormatInt(_plus, 2))
	//3
	//obj , err := NewJDataObj("./p2885362303.jpeg")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//obj.decodeImageData(context.Background())
	_chunk := make([]byte, 0)
	_chunk = append(_chunk, 246, 106, 84, 57, 151, 24)
	_buffer := allocArrStr(len(_chunk)*8, 1)
	for k, j := range _chunk {
		_byteStr := bin2Str(int64(rune(j)))
		s := 8 - len(_byteStr)
		if s > 0 {
			for i := 0; i < s; i++ {
				_byteStr = "0" + _byteStr
			}
		}
		for i := 0; i < 8; i++ {
			_buffer[k*8+i] = _byteStr[i]
		}
	}
	fmt.Println(string(_buffer), len(_buffer))
}
