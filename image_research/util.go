package image_research

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type HtDataArr struct {
	Key  string
	Item HfmTree
}

type HfmTree struct {
	Group int
	Value int
}

func SliceArr(data []byte, _n int) []byte {
	if _n+1 > len(data) {
		return []byte{}
	}
	return data[_n:]
}

func ReadInt16(chunk []byte, offset int) uint {
	return uint(chunk[offset])<<8 + uint(chunk[offset+1])<<0
}

func ExCount(str string) (zeroCount int, bitCount int) {
	noZero := false
	for _, v := range []byte(str) {
		if v == '1' {
			noZero = true
		}
		if !noZero {
			if v == '0' {
				zeroCount++
			}
		}
	}
	bitCount = len(str) - zeroCount
	return
}

func ReadInt32(chunk []byte, offset int) uint {
	return uint(chunk[offset])<<24 + uint(chunk[offset+1])<<16 + uint(chunk[offset+2])<<8 + uint(chunk[offset+3])<<0
}

func ReadInt8(chunk []byte, offset int) uint {
	return uint(chunk[offset]) << 0
}

//数字转8位二进制字符串
func NumberToString(num int64) string {
	res := strconv.FormatInt(num, 2)
	return fmt.Sprintf("%08s", res)
}

func NumberToStringByByte(num byte) string {
	if num >= 48 && num <= 57 {
		return NumberToString(int64(num))
	}
	return fmt.Sprintf("%08s", string(num))
}

func BufferToString(chunk []byte) string {
	return string(chunk)
}

func StrToBuffer(str string) []byte {
	return []byte(str)
}

func ReadBytes(chunk []byte, start, length uint) ([]byte, error) {
	end := start + length
	rest := make([]byte, 0)
	if int(end) > len(chunk) {
		return rest, errors.New("out of range ")
	}
	if len(chunk) >= int(end+1) {
		return chunk[int(start):int(end+1)], nil
	} else {
		return chunk[int(start):], nil
	}
}

func ReadBytesByStartAndEnd(chunk []byte, start, end uint) []byte {
	if start >= end {
		return []byte{}
	}
	if len(chunk) <= int(start) {
		return []byte{}
	}
	if len(chunk) > int(end) {
		return chunk[int(start):int(end)]
	} else {
		return chunk[int(start):]
	}
}

func RepeatString(str string, count int) string {
	return strings.Repeat(str, count)
}

func Shift(chunk []byte) (int, []byte) {
	if len(chunk) == 0 {
		return 0, chunk
	}
	return int(chunk[0]), chunk[1:]
}

//func tBin2sBin(b int) string {
//	base, _ := strconv.ParseInt(b, 2, 10)
//	return strconv.FormatInt(base, 16)
//}

func Str2Bin(_num string) (int64, error) {
	return strconv.ParseInt(_num, 2, 64)
}

func Bin2Str(_bin int64) string {
	return strconv.FormatInt(_bin, 2)
}

func FillBytes(_data []byte, fillByte byte) []byte {
	for k := range _data {
		_data[k] = fillByte
	}
	return _data
}

func Concat(_data1, _data2 []byte, length int) []byte {
	rows := _data1
	rows = append(rows, _data2...)
	if len(rows) >= length {
		return rows[0:length]
	}
	return rows
}

func AsciiToStr(_num int) string {
	return string(rune(_num))
}

func ByteToAscii(_one byte) int {
	return int(rune(_one))
}

func CreateHuffmanTree(chunk []byte, countArr []uint) (map[string]HfmTree, []HtDataArr, []byte) {
	ret := make(map[string]HfmTree)
	retArr := make([]HtDataArr, 0)
	last := ""
	for i := 0; i < len(countArr); i++ {
		_count := countArr[i]
		for j := 0; j < int(_count); j++ {
			if last == "" {
				last = RepeatString("0", i+1)
			} else {
				lastLength := len(last)
				_lastInt, _ := Str2Bin(last)
				last = Bin2Str(_lastInt + 1)

				if len(last) < lastLength {
					last = RepeatString("0", lastLength-len(last)) + last
				}
				if len(last) < i+1 {
					last = last + RepeatString("0", i+1-len(last))
				}
			}
			var _num int
			_num, chunk = Shift(chunk)
			_data := HfmTree{
				Group: len(last),
				Value: _num,
			}
			ret[last] = _data
			retArr = append(retArr, HtDataArr{
				Key:  last,
				Item: _data,
			})
		}
	}
	return ret, retArr, chunk
}
