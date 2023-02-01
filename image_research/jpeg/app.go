package jpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type colorComponentsArr struct {
	key  uint
	Item ColorComponentsItem
}

type HtArr struct {
	key  string
	item Ht
}

type JData struct {
	file                *os.File
	header              []byte
	index               int
	collectData         bool
	imageData           []JImageData
	accuracy            uint8 //样本精度，通常值为8，样本就是单个像素的颜色分量
	width               uint  //图片高度
	height              uint  //图片宽度
	colorComponentCount uint8 //颜色分量数，1 - 灰度图，3 - YCrCb/YIQ 彩色图，4 - CMYK 彩色图，通常值为3
	colorComponents     map[uint]ColorComponentsItem
	colorComponentsArr  []colorComponentsArr

	hMax      int64  //mcu宽度
	vMax      int64  //mcu高度
	readDuSeq []uint //mcu里数据单元的顺序，同样为Y- Y- Y- Y- Cb- Cr

	Ht        map[string]Ht
	colors    uint8 //1 - 灰度图，3 - YCrCb或YIQ，4 - CMYK
	colorInfo map[uint]ColorInfoItem

	qt map[uint]QtItem

	restartInterval uint //差分编码累计复位的间隔

	format          string //文件交换格式，通常是JFIF，JPEG File Interchange Format的缩写
	mainVersion     uint   //主版本号
	subVersion      uint   //次版本号
	unit            uint   //密度单位，0 - 无单位，1 - 点数/英寸，2 - 点数/厘米
	yPixel          uint   //垂直方向像素密度
	xPixel          uint   //水平方向像素密度
	thumbnailWidth  uint   //缩略图宽度
	thumbnailHeight uint   //缩略图高度

	thumbnail [][][]byte //缩略图

	comment []byte //注释

	pixels [][][]int
}

type QtItem struct {
	accuracy int64
	data     []byte
}

type ColorInfoItem struct {
	DCHTID int64
	ACHTID int64
}

type HtDataArr struct {
	key  string
	item HfmTree
}

type Ht struct {
	Type    string
	Data    map[string]HfmTree
	DataArr []HtDataArr
}

type HfmTree struct {
	Group int
	Value int
}

type ColorComponentsItem struct {
	X    int64
	Y    int64
	QtId uint
}

type JImageData struct {
	Type  []byte
	Chunk []byte
}

var startCode = []uint8{0xff, 0xd8}
var endCode = []uint8{0xff, 0xd9}

func NewFile(jpegFile string) (*JData, error) {
	fd, err := os.OpenFile(jpegFile, os.O_RDONLY, 0777)
	if err != nil {
		return nil , err
	}
	obj := &JData{file: fd, vMax: 1, hMax: 1,
		qt:              make(map[uint]QtItem),
		colorComponents: make(map[uint]ColorComponentsItem),
		Ht:              make(map[string]Ht),
		colorInfo:       make(map[uint]ColorInfoItem),
	}
	_, err = obj.decode(context.Background())
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (a *JData) exportHtml() error {
	canvas := fmt.Sprintf(`<canvas id="canvas" class="canvas" style="width: %dpx; height: %dpx;"></canvas>`, a.width, a.height)
	dataArr := make([]string, 0)
	for _, v := range a.pixels {
		for _, j := range v {
			dataArr = append(dataArr, fmt.Sprintf("%d,%d,%d,255", j[0], j[1], j[2]))
		}
	}
	data := strings.Join(dataArr, ",")
	script := fmt.Sprintf(`(function() {
    var canvas = document.getElementById('canvas');
    var ctx = canvas.getContext('2d');
    canvas.width = %d;
    canvas.height = %d;

    var pixels = Uint8ClampedArray.from([%s]);

    var imageData = new ImageData(pixels, %d, %d);
    ctx.putImageData(imageData, 0, 0, 0, 0, %d, %d);
})();`, a.width, a.height, data, a.width, a.height, a.width, a.height)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title></title>
    <style type="text/css">
        .cnt{
            display: flex;
            align-items: center;
        }

        .canvas{
            margin-left: 30px;
        }
    </style>
</head>
<body>
    <div class="cnt">
        %s
    </div>
    <script>
       	%s
    </script>
</body>
</html>`, canvas, script)
	res, err := os.Create("test.html")
	defer func() {
		_ = res.Close()
	}()
	if err != nil {
		return err
	}
	_, _ = res.WriteString(html)
	return nil
}

func NewJDataObj(jpegFile string) (*JData, error) {
	fd, err := os.OpenFile(jpegFile, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return &JData{file: fd, vMax: 1, hMax: 1}, nil
}

func (a *JData) decode(ctx context.Context) ([][][]int, error) {
	decodePixels := make([][][]int, 0)
	defer a.close(ctx)
	if err := a.decodeSOI(ctx); err != nil {
		return decodePixels, err
	}
	fileBytes, _ := ioutil.ReadAll(io.Reader(a.file))
	data := make([]byte, 0)
	for a.index < len(fileBytes) {
		oneByte, err := a.read(1)
		if err != nil {
			return decodePixels, err
		}
		_byte := oneByte[0]
		if _byte == 0xff {
			_next := fileBytes[a.index]
			if _next == 0x00 {
				_, _ = a.read(1)
				if a.collectData {
					data = append(data, _byte)
				}
			} else if _next == 0xff {
				continue
			} else if _next >= 0xd0 && _next <= 0xd7 {
				_, _ = a.read(1)
				a.imageData = append(a.imageData, JImageData{
					Type:  []byte{_next & 0x0f},
					Chunk: data,
				})
				data = []byte{}
			} else {
				if err := a.decodeMarkerSegment(ctx); err != nil {
					return decodePixels, err
				}
			}
		} else {
			if a.collectData {
				data = append(data, _byte)
			}
		}
	}
	if len(data) > 0 {
		a.imageData = append(a.imageData, JImageData{
			Type:  bytes.NewBufferString("end").Bytes(),
			Chunk: data,
		})
	}
	if len(a.pixels) > 0 {
		return a.pixels, nil
	}
	return a.decodeImageData(ctx)
}

func allocArrStr(_num int, fillStr byte) []byte {
	_re := make([]byte, 0)
	for i := 0; i < _num; i++ {
		_re = append(_re, fillStr)
	}
	return _re
}

func writeStr(replaceStr []byte, str []byte, offset int, length int) []byte {
	i := 0
	for k := range replaceStr {
		if k >= offset && k < offset+length {
			if len(str) > i {
				replaceStr[k] = str[i]
			}
			i++
		}
	}
	return replaceStr
}

type mcuArr struct {
	colorComponentId uint
	data             [][]int
}

func (a *JData) decodeImageData(ctx context.Context) ([][][]int, error) {
	mcus := make([][]mcuArr, 0)
	for _, v := range a.imageData {
		_chunk := v.Chunk
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
		readDuSeq := a.readDuSeq
		mcuDuCount := len(readDuSeq)
		mcu := make([]mcuArr, 0)
		rIndex := 0
		lastDc := make(map[uint]int)
		for len(_buffer) > 0 {
			colorComponentId := readDuSeq[rIndex]
			//哈夫曼解码
			cursor, output := a.decodeHuffman(ctx, _buffer, colorComponentId, lastDc[colorComponentId])

			lastDc[colorComponentId] = output[0]

			_buffer = sliceArr(_buffer, cursor)

			//反量化
			_output := a.quantify(output, colorComponentId, true)

			//zig zag反编码
			_output1 := zigZag(_output, true)

			//转矩阵
			_output2 := arrayToMatrix(_output1, 8, 8)

			//反dct编码
			_output3 := fastIDct(_output2)

			mcu = append(mcu, mcuArr{
				colorComponentId: colorComponentId,
				data:             _output3,
			})

			rIndex++
			if rIndex >= mcuDuCount {
				rIndex = 0

				mcus = append(mcus, mcu)
				mcu = []mcuArr{}
			}
		}
	}

	width := a.width
	height := a.height
	hPixels := a.hMax * 8
	vPixels := a.vMax * 8

	hNum := math.Ceil(float64(width) / float64(hPixels))
	vNum := math.Ceil(float64(height) / float64(vPixels))

	pixels := make([][][]int, width)
	for i := 0; i < int(width); i++ {
		pixels[i] = make([][]int, height)
	}
	for i := 0; i < int(hNum); i++ {
		for j := 0; j < int(vNum); j++ {
			mcuPixels := a.getMcuPixels(mcus[j*int(hNum)+i])

			offsetX := i * int(hPixels)
			offsetY := j * int(vPixels)
			for x := 0; x < int(hPixels); x++ {
				for y := 0; y < int(vPixels); y++ {
					insertX := x + offsetX
					insertY := y + offsetY

					if insertX < int(width) && insertY < int(height) {
						pixels[insertX][insertY] = ycrcb2rgb(mcuPixels[x][y])
					}
				}
			}
		}
	}
	a.pixels = pixels
	return pixels, nil
}

var C1 = math.Cos(math.Pi / 16)
var C2 = math.Cos(math.Pi / 8)
var C3 = math.Cos(3 * math.Pi / 16)
var C4 = math.Cos(math.Pi / 4)
var C5 = math.Cos(5 * math.Pi / 16)
var C6 = math.Cos(3 * math.Pi / 8)
var C7 = math.Cos(7 * math.Pi / 16)

// 亮度量化表
var B_QUANTIZATION_TABLE = []int{
	16, 12, 14, 14, 18, 24, 49, 72,
	11, 12, 13, 17, 22, 35, 64, 92,
	10, 14, 16, 22, 37, 55, 78, 95,
	16, 19, 24, 29, 56, 64, 87, 98,
	24, 26, 40, 51, 68, 81, 103, 112,
	40, 58, 57, 87, 109, 104, 121, 100,
	51, 60, 69, 80, 103, 113, 120, 103,
	61, 55, 56, 62, 77, 92, 101, 99,
}

// 色度量化表
var C_QUANTIZATION_TABLE = []int{
	17, 18, 24, 47, 99, 99, 99, 99,
	18, 21, 26, 66, 99, 99, 99, 99,
	24, 26, 56, 99, 99, 99, 99, 99,
	47, 66, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
	99, 99, 99, 99, 99, 99, 99, 99,
}

func (a *JData) quantify(input []int, colorComponentId uint, isReverse bool) []int {
	output := make([]int, 64)
	if !isReverse {
		qt := C_QUANTIZATION_TABLE
		if colorComponentId == 1 {
			qt = B_QUANTIZATION_TABLE
		}
		for i := 0; i < 64; i++ {
			output[i] = int(math.Round(float64(input[i] / qt[i])))
		}
	} else {
		colorComponent := a.colorComponents[colorComponentId]
		qt := a.qt[colorComponent.QtId].data
		//fmt.Println("colorComponent", colorComponent, "qt", qt, "colorComponent.QtId", colorComponent.QtId)
		//fmt.Println(input, len(input), output)
		for i := 0; i < 64; i++ {
			output[i] = int(math.Round(float64(input[i] / int(qt[i]))))
		}
	}

	return output
}

func arrayToMatrix(input []int, w, h int) [][]int {
	output := make([][]int, w)

	for i := 0; i < w; i++ {
		output[i] = make([]int, h)
		for j := 0; j < h; j++ {
			output[i][j] = input[j*w+i]
		}
	}
	return output
}

var ZIG_ZAG = []int{
	0, 1, 5, 6, 14, 15, 27, 28,
	2, 4, 7, 13, 16, 26, 29, 42,
	3, 8, 12, 17, 25, 30, 41, 43,
	9, 11, 18, 24, 31, 40, 44, 53,
	10, 19, 23, 32, 39, 45, 52, 54,
	20, 22, 33, 38, 46, 51, 55, 60,
	21, 34, 37, 47, 50, 56, 59, 61,
	35, 36, 48, 49, 57, 58, 62, 63,
}

func zigZag(input []int, isReverse bool) []int {
	output := make([]int, 64)
	for i := 0; i < 64; i++ {
		if len(input) >= 64 {
			if isReverse {
				output[i] = input[ZIG_ZAG[i]]
			} else {
				output[ZIG_ZAG[i]] = input[i]
			}
		}
	}

	return output
}

//快速逆dct变换
func fastIDct(input [][]int) [][]int {
	output := make([][]int, 8)

	for i := 0; i < 8; i++ {
		output[i] = make([]int, 8)
		for j := 0; j < 8; j++ {
			output[i][j] = int(input[i][j])
		}
	}
	for i := 0; i < 8; i++ {
		_temp0 := (float64(output[0][i])*C4 + float64(output[2][i])*C2 + float64(output[4][i])*C4 + float64(output[6][i])*C6) / 2
		_temp1 := (float64(output[0][i])*C4 + float64(output[2][i])*C6 - float64(output[4][i])*C4 - float64(output[6][i])*C2) / 2
		_temp2 := (float64(output[0][i])*C4 - float64(output[2][i])*C6 - float64(output[4][i])*C4 + float64(output[6][i])*C2) / 2
		_temp3 := (float64(output[0][i])*C4 - float64(output[2][i])*C2 + float64(output[4][i])*C4 - float64(output[6][i])*C6) / 2
		_temp4 := (float64(output[1][i])*C7 - float64(output[3][i])*C5 + float64(output[5][i])*C3 - float64(output[7][i])*C1) / 2
		_temp5 := (float64(output[1][i])*C5 - float64(output[3][i])*C1 + float64(output[5][i])*C7 + float64(output[7][i])*C3) / 2
		_temp6 := (float64(output[1][i])*C3 - float64(output[3][i])*C7 - float64(output[5][i])*C1 - float64(output[7][i])*C5) / 2
		_temp7 := (float64(output[1][i])*C1 + float64(output[3][i])*C3 + float64(output[5][i])*C5 + float64(output[7][i])*C7) / 2

		output[0][i] = int(_temp0 + _temp7)
		output[1][i] = int(_temp1 + _temp6)
		output[2][i] = int(_temp2 + _temp5)
		output[3][i] = int(_temp3 + _temp4)
		output[4][i] = int(_temp3 + _temp4)
		output[5][i] = int(_temp2 + _temp5)
		output[6][i] = int(_temp1 + _temp6)
		output[7][i] = int(_temp0 + _temp7)
	}

	for i := 0; i < 8; i++ {
		_temp0 := (float64(output[i][0])*C4 + float64(output[i][2])*C2 + float64(output[i][4])*C4 + float64(output[i][6])*C6) / 2
		_temp1 := (float64(output[i][0])*C4 + float64(output[i][2])*C6 - float64(output[i][4])*C4 - float64(output[i][6])*C2) / 2
		_temp2 := (float64(output[i][0])*C4 - float64(output[i][2])*C6 - float64(output[i][4])*C4 + float64(output[i][6])*C2) / 2
		_temp3 := (float64(output[i][0])*C4 - float64(output[i][2])*C2 + float64(output[i][4])*C4 - float64(output[i][6])*C6) / 2
		_temp4 := (float64(output[i][0])*C7 - float64(output[i][3])*C5 + float64(output[i][5])*C3 - float64(output[i][7])*C1) / 2
		_temp5 := (float64(output[i][0])*C5 - float64(output[i][3])*C1 + float64(output[i][5])*C7 + float64(output[i][7])*C3) / 2
		_temp6 := (float64(output[i][0])*C3 - float64(output[i][3])*C7 - float64(output[i][5])*C1 - float64(output[i][7])*C5) / 2
		_temp7 := (float64(output[i][0])*C1 + float64(output[i][3])*C3 + float64(output[i][5])*C5 + float64(output[i][7])*C7) / 2

		output[i][0] = int(_temp0 + _temp7)
		output[i][1] = int(_temp1 + _temp6)
		output[i][2] = int(_temp2 + _temp5)
		output[i][3] = int(_temp3 + _temp4)
		output[i][4] = int(_temp3 + _temp4)
		output[i][5] = int(_temp2 + _temp5)
		output[i][6] = int(_temp1 + _temp6)
		output[i][7] = int(_temp0 + _temp7)
	}
	return output
}

type splitArr struct {
	Dus   [][][]int
	Index int
	X     int64
	Y     int64
}

func (a *JData) getMcuPixels(mcu []mcuArr) [][][]int {
	hmax := a.hMax
	vmax := a.vMax

	hPixels := hmax * 8
	vPixels := vmax * 8

	output := make([][][]int, int(a.width))
	for i := 0; i < int(hPixels); i++ {
		output[i] = make([][]int, int(a.height))
		for j := 0; j < int(vPixels); j++ {
			output[i][j] = make([]int, 0)
		}
	}

	tempArr := make([]*splitArr, 0)
	tempArr = append(tempArr, &splitArr{
		Dus:   make([][][]int, 0),
		Index: 0,
	}, &splitArr{
		Dus:   make([][][]int, 0),
		Index: 2,
	}, &splitArr{
		Dus:   make([][][]int, 0),
		Index: 1,
	}) //0-Y, 1-Cb, 2-Cr

	for _, v := range mcu {
		colorComponentId := v.colorComponentId

		_temp := tempArr[colorComponentId-1]
		_temp.X = a.colorComponents[colorComponentId].X
		_temp.Y = a.colorComponents[colorComponentId].Y
		_temp.Dus = append(_temp.Dus, v.data)
	}

	for i := 0; i < int(hPixels); i++ {
		for j := 0; j < int(vPixels); j++ {
			output[i][j] = make([]int, 3)
		}
	}

	for _, v := range tempArr {
		dus := v.Dus
		xRange := hmax / v.X //采样块宽度
		yRange := vmax / v.Y //采样块高度
		data := make([][]int, a.width)
		for i := 0; i < len(data); i++ {
			data[i] = make([]int, a.height)
		}

		for h := 0; h < int(v.X); h++ {
			for x := 0; x < int(v.Y); x++ {
				du := dus[x*h+h]
				for i := 0; i < 8; i++ {
					insertX := (int(i) + int(h)*8) * int(xRange) //采样点少的情况下需要步偏移

					for j := 0; j < 8; j++ {
						insertY := (j + x*8) * int(yRange)
						value := du[i][j]
						for rx := 0; rx < int(xRange); rx++ {
							for ry := 0; ry < int(yRange); ry++ {
								data[insertX+rx][insertY] = value
							}
						}
					}
				}
			}
		}

		index := v.Index
		for i := 0; i < int(hPixels); i++ {
			for j := 0; j < int(vPixels); j++ {
				output[i][j][index] = data[i][j] + 128
			}
		}
	}

	return output
}

/*
 * YCrCb 转 RGB
 * [ r ]   [ 1  1.402     0        ][ Y        ]
 * [ g ] = [ 1  -0.71414  -0.34414 ][ Cr - 128 ]
 * [ b ]   [ 1  0         1.772    ][ Cb - 128 ]
 */
func ycrcb2rgb(YCrCb []int) []int {
	r := float64(YCrCb[0]) + 1.402*(float64(YCrCb[1])-128)
	g := float64(YCrCb[0]) - 0.34414*(float64(YCrCb[2])-128) - 0.71414*(float64(YCrCb[1])-128)
	b := float64(YCrCb[0]) + 1.772*(float64(YCrCb[2])-128)

	_r := 0
	_g := 0
	_b := 0
	if r > 255 {
		_r = 255
	} else {
		if r > 0 {
			_r = int(r)
		}
	}
	if g > 255 {
		_g = 255
	} else {
		if g > 0 {
			_g = int(g)
		}
	}
	if b > 255 {
		_b = 255
	} else {
		if b > 0 {
			_b = int(b)
		}
	}
	return []int{
		_r,
		_g,
		_b,
	}
}

/**
 * RGB 转 YCrCb
 *
 * [ Y  ]   [ 0.299    0.587    0.244   ][ r ]   [ 0   ]
 * [ Cr ] = [ 0.5      -0.4187  -0.0813 ][ g ] + [ 128 ]
 * [ Cb ]   [ -0.1687  -0.3313  0.5     ][ b ]   [ 128 ]
 */
func rgb2ycrcb(rgb []int) []float64 {
	return []float64{
		0.299*float64(rgb[0]) + 0.587*float64(rgb[1]) + 0.114*float64(rgb[2]),
		0.5*float64(rgb[0]) - 0.4187*float64(rgb[1]) - 0.0813*float64(rgb[2]) + 128,
		-0.1687*float64(rgb[0]) - 0.3313*float64(rgb[1]) + 0.5*float64(rgb[2]) + 128,
	}
}

func keys(list []HtDataArr) []string {
	_list := make([]string, 0)
	for _, v := range list {
		_list = append(_list, v.key)
	}
	return _list
}

func (a *JData) decodeHuffman(ctx context.Context, input []byte, colorComponentId uint, lastDc int) (int, []int) {
	cursor := 0
	output := make([]int, 0)
	for i := 0; i < 64; i++ {
		_r := a.colorInfo[colorComponentId].DCHTID
		var _d int
		if i == 0 {
			_d = 0
			_r = a.colorInfo[colorComponentId].ACHTID
		} else {
			_d = 1
		}
		_ht := a.Ht[fmt.Sprintf("%d-%d", _d, _r)].DataArr
		_htMap := a.Ht[fmt.Sprintf("%d-%d", _d, _r)].Data
		_keys := keys(_ht)
		//fmt.Println("_keys", _keys)
		//fmt.Println("_ht", _ht)
		var value int
		for _, v := range _keys {
			length := len(v)
			//fmt.Println("v---", v, len(v))
			var subBuffer []byte
			if len(input) >= cursor+length {
				subBuffer = readBytesByStartAndEnd(input, uint(cursor), uint(length))
			} else {
				subBuffer = sliceArr(input, cursor)
				tempSubBuffer := make([]byte, length-len(subBuffer))
				tempSubBuffer = fillBytes(tempSubBuffer, '1')
				subBuffer = concat(subBuffer, tempSubBuffer, length)
			}
			if string(subBuffer) == v {
				cursor += length
				value = _htMap[v].Value
				break
			}
		}
		//fmt.Println("cur ", cursor)

		var bitCount int
		var bitData int
		// dc、ac数据的解法比较特殊，类似如下
		// 位数                值
		// 0                   0
		// 1                 -1, 1
		// 2             -3, -2, 2, 3
		// 3        -7, ..., -4, 4, ..., 7
		// 4       -15, ..., -8, 8, ..., 15
		if i == 0 {
			//取DC值
			bitCount = value

			if bitCount == 0 {
				bitData = 0
			} else {
				_bitData := readBytesByStartAndEnd(input, uint(cursor), uint(bitCount+cursor))
				_bitData1, _ := str2Bin(string(_bitData))
				half := math.Pow(2, float64(bitCount-1))
				if float64(_bitData1) >= half {
					bitData = int(_bitData1)
				} else {
					bitData = int(float64(_bitData1) - half*2 + 1)
				}
			}
			bitData += lastDc
		} else {
			//取AC值
			bitString := numberToString(int64(value))
			zeroCount, _ := strconv.ParseInt(bitString[0:4], 10, 2) //数据前0的个数
			_bitCount, _ := strconv.ParseInt(bitString[4:8], 10, 2) //数据的位数
			bitCount = int(_bitCount)
			//fmt.Println("bitString", bitString,
			//	"zeroCount", zeroCount,
			//	"_bitCount", _bitCount,
			//	"bitCount", bitCount,
			//)

			if zeroCount == 0 && bitCount == 0 {
				//解析到 (0, 0) ，表示到达EOB，意味着后面的都是0
				for len(output) < 64 {
					output = append(output, 0)
				}
				break
			} else {
				if bitCount == 0 {
					bitData = 0
				} else {
					_bitData := readBytesByStartAndEnd(input, uint(cursor), uint(bitCount+cursor))
					_bitData1, _ := str2Bin(string(_bitData))
					half := math.Pow(2, float64(bitCount-1))
					if float64(_bitData1) >= half {
						bitData = int(_bitData1)
					} else {
						bitData = int(float64(_bitData1) - half*2 + 1)
					}
				}
				for j := 0; j < int(zeroCount); j++ {
					output = append(output, 0)
					i++
				}
			}
		}
		//fmt.Println("bc---", bitCount, cursor)
		output = append(output, bitData)
		cursor += bitCount
	}
	//fmt.Println("res ---", cursor, output)
	return cursor, output
}

func (a *JData) decodeMarkerSegment(ctx context.Context) error {
	markerB, err := a.read(1)
	if err != nil {
		return err
	}
	marker := markerB[0]
	if marker == 0xd9 {
		return nil
	}
	length, _ := a.read(2)
	lenInt := length[0] + length[1]
	chunkData, _ := a.read(int64(lenInt) + -2)
	switch marker {
	case 0xc0:
		if err := a.decodeSOf0(ctx, chunkData); err != nil {
			return err
		}
	case 0xc4:
		if err := a.decodeDHT(ctx, chunkData); err != nil {
			return err
		}
	case 0xda:
		a.decodeSOS(ctx, chunkData)
	case 0xdb:
		if err := a.decodeDQT(ctx, chunkData); err != nil {
			return err
		}
	case 0xdd:
		a.decodeDRI(ctx, chunkData)
	case 0xe0:
		if err := a.decodeAPP0(ctx, chunkData); err != nil {
			return err
		}
	case 0xfe:
		a.decodeCOM(ctx, chunkData)
	}
	return nil
}

func (a *JData) decodeSOf0(ctx context.Context, chunkData []byte) error {
	a.accuracy = chunkData[0]
	a.height = readInt16(chunkData, 1)
	a.width = readInt16(chunkData, 3)
	a.colorComponentCount = chunkData[5]
	if a.colorComponentCount != 3 {
		return errors.New("仅支持YCrCb彩色图")
	}
	chunk := sliceArr(chunkData, 6)
	a.colorComponents = make(map[uint]ColorComponentsItem)
	for i := 0; i < int(a.colorComponentCount); i++ {
		colorComponentId := readInt8(chunk, 0)
		packField := readInt8(chunk, 1)
		_pF := numberToString(int64(packField))
		_x, _ := strconv.ParseInt(_pF[0:4], 2, 10)
		_y, _ := strconv.ParseInt(_pF[4:8], 2, 10)
		qrId := readInt8(chunk, 2)
		_item := ColorComponentsItem{
			X:    _x,
			Y:    _y,
			QtId: qrId,
		}
		a.colorComponents[colorComponentId] = _item
		a.colorComponentsArr = append(a.colorComponentsArr, colorComponentsArr{
			key:  colorComponentId,
			Item: _item,
		})
		chunk = sliceArr(chunk, 3)
	}
	for _, v := range a.colorComponentsArr {
		for i := 0; i < int(v.Item.X)*int(v.Item.Y); i++ {
			a.readDuSeq = append(a.readDuSeq, v.key)
		}
		if v.Item.X > a.hMax {
			a.hMax = v.Item.X
		}
		if v.Item.Y > a.vMax {
			a.vMax = v.Item.Y
		}
	}
	return nil
}

func (a *JData) decodeDHT(ctx context.Context, chunkData []byte) error {
	for len(chunkData) > 0 {
		packedField := readInt8(chunkData, 0)
		_pf := numberToString(int64(packedField))
		_type := _pf[3] //哈夫曼类型 0-DC直流 1-AC交流表
		_id := _pf[7]   //哈夫曼表id
		var length uint = 0
		countArr := make([]uint, 0)
		for i := 0; i < 16; i++ {
			_count := readInt8(chunkData, i+1)
			countArr = append(countArr, _count)
			length += _count
		}
		data, err := readBytes(chunkData, 17, length)
		if err != nil {
			return err
		}
		_data, _dataArr, _ := createHuffmanTree(data, countArr)
		_ht := Ht{
			Type:    asciiToStr(int(_type)),
			Data:    _data,
			DataArr: _dataArr,
		}
		_key := fmt.Sprintf("%s-%s", asciiToStr(int(_type)), asciiToStr(int(_id)))
		a.Ht[_key] = _ht
		chunkData = sliceArr(chunkData, int(length)+17)
	}
	return nil
}

func (a *JData) decodeSOS(ctx context.Context, chunkData []byte) {
	a.colors = uint8(readInt8(chunkData, 0))
	chunkData = sliceArr(chunkData, 1)
	for i := 0; i < int(a.colors); i++ {
		colorId := readInt8(chunkData, 0) //1 - Y，2 - Cb，3 - Cr，4 - I，5 - Q
		packedField := readInt8(chunkData, 1)
		_pf := numberToString(int64(packedField))

		_dchtId, _ := strconv.ParseInt(_pf[0:4], 2, 10) //DC哈夫曼表id
		_achtId, _ := strconv.ParseInt(_pf[4:8], 2, 10) //AC哈夫曼表id

		a.colorInfo[colorId] = ColorInfoItem{
			DCHTID: _dchtId,
			ACHTID: _achtId,
		}
		chunkData = sliceArr(chunkData, 2)
	}
	a.collectData = true //准备开始收集数据
}

//解码DQT，Define Quantization Table，定义量化表
func (a *JData) decodeDQT(ctx context.Context, chunkData []byte) error {
	for len(chunkData) > 0 {
		packedField := readInt8(chunkData, 0)
		_pf := numberToString(int64(packedField))

		accuracy, _ := strconv.ParseInt(_pf[0:4], 2, 10) //量化表精度，0 - 1字节，1 - 2字节
		id, _ := strconv.ParseInt(_pf[4:8], 2, 10)       //量化表id，取值范围为0 - 3，所以最多可有4个量化表

		length := 64 * (accuracy + 1)
		data, err := readBytes(chunkData, 1, uint(length))
		if err != nil {
			return err
		}
		a.qt[uint(id)] = QtItem{
			accuracy: accuracy,
			data:     data,
		}

		chunkData = sliceArr(chunkData, len(chunkData)+1)
	}
	return nil
}

//解码DRI，Define Restart Interval，定义差分编码累计复位的间隔
func (a *JData) decodeDRI(ctx context.Context, chunkData []byte) {
	a.restartInterval = readInt16(chunkData, 0)
}

//Application，应用程序保留标记0
func (a *JData) decodeAPP0(ctx context.Context, chunkData []byte) error {
	_formatB, err := readBytes(chunkData, 0, 5)
	if err != nil {
		return err
	}
	a.format = bufferToString(_formatB)
	a.mainVersion = readInt8(chunkData, 5)
	a.subVersion = readInt8(chunkData, 6)

	a.unit = readInt8(chunkData, 7)
	a.xPixel = readInt16(chunkData, 8)
	a.yPixel = readInt16(chunkData, 10)

	a.thumbnailWidth = readInt8(chunkData, 12)
	a.thumbnailHeight = readInt8(chunkData, 13)

	if a.thumbnailWidth > 0 && a.thumbnailHeight > 0 && len(chunkData) > 14 {
		thumbnail := sliceArr(chunkData, 14)
		data := make([][][]byte, 0)
		for i := 0; i < int(a.thumbnailWidth); i++ {
			data[i] = make([][]byte, 0)
			for j := 0; j < int(a.thumbnailHeight); j++ {
				index := j*int(a.thumbnailWidth) + i

				data[i][j] = []byte{
					thumbnail[index],
					thumbnail[index+1],
					thumbnail[index+2],
				}
			}
		}
		a.thumbnail = data
	}

	return nil
}

func (a *JData) decodeCOM(ctx context.Context, chunkData []byte) {
	a.comment = chunkData
}

func (a *JData) read(_n int64) ([]byte, error) {
	_bytes := make([]byte, _n)
	_, err := a.file.ReadAt(_bytes, int64(a.index))
	if err != nil {
		return _bytes, err
	}
	a.index += int(_n)
	return _bytes, nil
}

func (a *JData) decodeSOI(ctx context.Context) error {
	_bytes, err := a.read(2)
	if err != nil {
		return err
	}
	if _bytes[0] != startCode[0] || _bytes[1] != startCode[1] {
		return errors.New("not invalid jpeg format")
	}
	a.index = 2
	a.header = _bytes
	return nil
}

func (a *JData) close(ctx context.Context) {
	_ = a.file.Close()
}
