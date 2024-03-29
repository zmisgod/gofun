package jpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/zmisgod/gofun/image_research"
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
	Data    map[string]image_research.HfmTree
	DataArr []image_research.HtDataArr
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
		return nil, err
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
	for i := 0; i < int(a.height); i++ {
		for j := 0; j < int(a.width); j++ {
			_item := a.pixels[j][i]
			dataArr = append(dataArr, fmt.Sprintf("%d,%d,%d,255", _item[0], _item[1], _item[2]))
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

/*
--Other markers--
X’FFD8’ |  SOI* | Start of image
X’FFD9’ |  EOI* | End of image
X’FFDA’ |  SOS  | Start of scan
X’FFDB’ |  DQT  | Define quantization table(s)
X’FFDC’ |  DNL  | Define number of lines
X’FFDD’ |  DRI  | Define restart interval
X’FFDE’ |  DHP  | Define hierarchical progression
X’FFDF’ |  EXP  | Expand reference component(s)
X’FFE0’ through X’FFEF’ | APPn | Reserved for application segments
X’FFF0’ through X’FFFD’ | JPGn | Reserved for JPEG extensions
X’FFFE’ |  COM  | Comment

--Start Of Frame markers, non-differential, Huffman coding--
X’FFC0’｜SOF0｜Baseline DCT
X’FFC1’｜SOF1｜Extended sequential DCT
X’FFC2’｜SOF2｜Progressive DCT
X’FFC3’｜SOF3｜Lossless (sequential)

--Start Of Frame markers, differential, Huffman coding--
X’FFC5’｜SOF5｜Differential sequential DCT
X’FFC6’｜SOF6｜Differential progressive DCT
X’FFC7’｜SOF7｜Differential lossless (sequential)

--Start Of Frame markers, non-differential, arithmetic coding--
X’FFC8’｜JPG｜Reserved for JPEG extensions
X’FFC9’｜SOF9｜Extended sequential DCT
X’FFCA’｜SOF10｜Progressive DCT
X’FFCB’｜SOF11｜Lossless (sequential)

--Start Of Frame markers, differential, arithmetic coding--
X’FFCD’｜SOF13｜Differential sequential DCT
X’FFCE’｜SOF14｜Differential progressive DCT
X’FFCF’｜SOF15｜Differential lossless (sequential)

--Huffman table specification--
X’FFC4’｜DHT｜Define Huffman table(s)

--Arithmetic coding conditioning specification--
X’FFCC’｜ DAC ｜Define arithmetic coding conditioning(s)

--Restart interval termination--
X’FFD0’ through X’FFD7’｜ RSTm*｜ Restart with modulo 8 count “m”

--Reserved markers--
X’FF01’｜TEM*｜For temporary private use in arithmetic coding
X’FF02’ through X’FFBF’｜RES｜Reserved
 */
func (a *JData) decode(ctx context.Context) (decodePixels [][][]int, err error) {
	defer a.close(ctx)
	//检查前2个字节是否为0xff, 0xd8，这是正确的jpeg格式前缀
	err = a.decodeSOI(ctx)
	if err != nil {
		return
	}
	fileBytes, _ := ioutil.ReadAll(io.Reader(a.file))
	data := make([]byte, 0)
	for a.index < len(fileBytes) {
		//读一个字节，判断该字节是否为0xff 段标识符
		oneByte, _ := a.read(1)
		_byte := oneByte[0]
		if _byte == 0xff {
			//继续往下面读取数据
			_next := fileBytes[a.index]
			if _next == 0x00 {
				//图像里的一部分
				_, _ = a.read(1)//跳过next
				if a.collectData {
					//保存数据
					data = append(data, _byte)
				}
			} else if _next == 0xff {
				//跳过分隔符
				continue
			} else if _next >= 0xd0 && _next <= 0xd7 {
				//Restart interval termination
				//遇到RSTn标记
				_, _ = a.read(1)//跳过next
				a.imageData = append(a.imageData, JImageData{
					Type:  []byte{_next & 0x0f},
					Chunk: data,
				})
				data = []byte{}
			} else {
				//解析其他标记
				if err := a.decodeMarkerSegment(ctx); err != nil {
					return decodePixels, err
				}
			}
		} else {
			//不是分隔符，就需要收集图片数据
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
	data             [][]float64
}

func (a *JData) decodeImageData(ctx context.Context) ([][][]int, error) {
	mcus := make([][]mcuArr, 0)
	for _, v := range a.imageData {
		_chunk := v.Chunk
		_buffer := allocArrStr(len(_chunk)*8, 1)
		for k, j := range _chunk {
			_byteStr := image_research.Bin2Str(int64(rune(j)))
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

			_buffer = image_research.SliceArr(_buffer, cursor)

			//反量化
			_output := a.quantify(output, colorComponentId, true)
			//fmt.Println("output", _output)

			//zig zag反编码
			_output1 := zigZag(_output, true)

			//转矩阵
			_output2 := arrayToMatrix(_output1, 8, 8)

			//fmt.Println(_output2)

			//反dct编码
			_output3 := fastIDct(_output2)
			//for _, _x := range _output3 {
			//	fmt.Println("---", _x)
			//}

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
	//fmt.Println(hNum, vNum, hPixels, vPixels, width, height)

	//for _, _x := range mcus {
	//	for _, _j := range _x {
	//		for _, _k := range _j.data {
	//			fmt.Println("---", _k)
	//		}
	//	}
	//}

	for i := 0; i < int(hNum); i++ {
		for j := 0; j < int(vNum); j++ {
			mcuPixels := a.getMcuPixels(mcus[j*int(hNum)+i], 0, 0)
			//fmt.Println("---", j*int(hNum)+i, i, j)

			offsetX := i * int(hPixels)
			offsetY := j * int(vPixels)
			for x := 0; x < int(hPixels); x++ {
				for y := 0; y < int(vPixels); y++ {
					insertX := x + offsetX
					insertY := y + offsetY

					if insertX < int(width) && insertY < int(height) {
						_a := ycrcb2rgb(mcuPixels[x][y])
						pixels[insertX][insertY] = _a
						fmt.Println(x, y, _a, mcuPixels[x][y])
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

func (a *JData) quantify(input []int, colorComponentId uint, isReverse bool) []float64 {
	output := make([]float64, 64)
	if !isReverse {
		qt := C_QUANTIZATION_TABLE
		if colorComponentId == 1 {
			qt = B_QUANTIZATION_TABLE
		}
		for i := 0; i < 64; i++ {
			output[i] = math.Round(float64(input[i] / qt[i]))
		}
	} else {
		colorComponent := a.colorComponents[colorComponentId]
		qt := a.qt[colorComponent.QtId].data
		//fmt.Println("colorComponent", colorComponent, "qt", qt, "colorComponent.QtId", colorComponent.QtId)
		//fmt.Println(input, len(input), output)
		for i := 0; i < 64; i++ {
			output[i] = math.Round(float64(input[i] * int(qt[i])))
		}
	}

	return output
}

func arrayToMatrix(input []float64, w, h int) [][]float64 {
	output := make([][]float64, w)

	for i := 0; i < w; i++ {
		output[i] = make([]float64, h)
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

func zigZag(input []float64, isReverse bool) []float64 {
	output := make([]float64, 64)
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
func fastIDct(input [][]float64) [][]float64 {
	output := make([][]float64, 8)
	for i := 0; i < 8; i++ {
		output[i] = make([]float64, 8)
		for j := 0; j < 8; j++ {
			output[i][j] = input[i][j]
		}
	}
	for i := 0; i < 8; i++ {
		_temp0 := (float64(output[0][i])*C4 + float64(output[2][i])*C2 + float64(output[4][i])*C4 + float64(output[6][i])*C6) / 2.0
		_temp1 := (float64(output[0][i])*C4 + float64(output[2][i])*C6 - float64(output[4][i])*C4 - float64(output[6][i])*C2) / 2.0
		_temp2 := (float64(output[0][i])*C4 - float64(output[2][i])*C6 - float64(output[4][i])*C4 + float64(output[6][i])*C2) / 2.0
		_temp3 := (float64(output[0][i])*C4 - float64(output[2][i])*C2 + float64(output[4][i])*C4 - float64(output[6][i])*C6) / 2.0
		_temp4 := (float64(output[1][i])*C7 - float64(output[3][i])*C5 + float64(output[5][i])*C3 - float64(output[7][i])*C1) / 2.0
		_temp5 := (float64(output[1][i])*C5 - float64(output[3][i])*C1 + float64(output[5][i])*C7 + float64(output[7][i])*C3) / 2.0
		_temp6 := (float64(output[1][i])*C3 - float64(output[3][i])*C7 - float64(output[5][i])*C1 - float64(output[7][i])*C5) / 2.0
		_temp7 := (float64(output[1][i])*C1 + float64(output[3][i])*C3 + float64(output[5][i])*C5 + float64(output[7][i])*C7) / 2.0

		output[0][i] = _temp0 + _temp7
		output[1][i] = _temp1 + _temp6
		output[2][i] = _temp2 + _temp5
		output[3][i] = _temp3 + _temp4
		output[4][i] = _temp3 - _temp4
		output[5][i] = _temp2 - _temp5
		output[6][i] = _temp1 - _temp6
		output[7][i] = _temp0 - _temp7
	}

	for i := 0; i < 8; i++ {
		_temp0 := (float64(output[i][0])*C4 + float64(output[i][2])*C2 + float64(output[i][4])*C4 + float64(output[i][6])*C6) / 2.0
		_temp1 := (float64(output[i][0])*C4 + float64(output[i][2])*C6 - float64(output[i][4])*C4 - float64(output[i][6])*C2) / 2.0
		_temp2 := (float64(output[i][0])*C4 - float64(output[i][2])*C6 - float64(output[i][4])*C4 + float64(output[i][6])*C2) / 2.0
		_temp3 := (float64(output[i][0])*C4 - float64(output[i][2])*C2 + float64(output[i][4])*C4 - float64(output[i][6])*C6) / 2.0
		_temp4 := (float64(output[i][1])*C7 - float64(output[i][3])*C5 + float64(output[i][5])*C3 - float64(output[i][7])*C1) / 2.0
		_temp5 := (float64(output[i][1])*C5 - float64(output[i][3])*C1 + float64(output[i][5])*C7 + float64(output[i][7])*C3) / 2.0
		_temp6 := (float64(output[i][1])*C3 - float64(output[i][3])*C7 - float64(output[i][5])*C1 - float64(output[i][7])*C5) / 2.0
		_temp7 := (float64(output[i][1])*C1 + float64(output[i][3])*C3 + float64(output[i][5])*C5 + float64(output[i][7])*C7) / 2.0

		output[i][0] = _temp0 + _temp7
		output[i][1] = _temp1 + _temp6
		output[i][2] = _temp2 + _temp5
		output[i][3] = _temp3 + _temp4
		output[i][4] = _temp3 - _temp4
		output[i][5] = _temp2 - _temp5
		output[i][6] = _temp1 - _temp6
		output[i][7] = _temp0 - _temp7
	}
	return output
}

type splitArr struct {
	Dus   [][][]float64
	Index int
	X     int64
	Y     int64
}

func (a *JData) getMcuPixels(mcu []mcuArr, _i, _j int) [][][]float64 {
	hmax := a.hMax
	vmax := a.vMax

	hPixels := hmax * 8
	vPixels := vmax * 8

	output := make([][][]float64, int(hPixels))
	for i := 0; i < int(hPixels); i++ {
		output[i] = make([][]float64, int(vPixels))
		for j := 0; j < int(vPixels); j++ {
			output[i][j] = make([]float64, 0)
		}
	}
	debug := false
	if _i == 0 && _j == 1 {
		debug = true
	}
	if debug {
		//for _, v := range mcu {
		//	fmt.Println("---aaa", v)
		//}
	}

	tempArr := make([]*splitArr, 0)
	tempArr = append(tempArr, &splitArr{
		Dus:   make([][][]float64, 0),
		Index: 0,
	}, &splitArr{
		Dus:   make([][][]float64, 0),
		Index: 2,
	}, &splitArr{
		Dus:   make([][][]float64, 0),
		Index: 1,
	}) //0-Y, 1-Cb, 2-Cr

	for _, v := range mcu {
		colorComponentId := v.colorComponentId

		_temp := tempArr[colorComponentId-1]
		_temp.X = a.colorComponents[colorComponentId].X
		_temp.Y = a.colorComponents[colorComponentId].Y
		tempArr[colorComponentId-1].Dus = append(tempArr[colorComponentId-1].Dus, v.data)
	}
	if debug {
		for _, v := range tempArr {
			fmt.Println("---bbb", v)
		}
	}

	for i := 0; i < int(hPixels); i++ {
		for j := 0; j < int(vPixels); j++ {
			output[i][j] = make([]float64, 3)
		}
	}

	for _, v := range tempArr {
		dus := v.Dus
		xRange := hmax / v.X //采样块宽度
		yRange := vmax / v.Y //采样块高度
		data := make([][]float64, hPixels)
		for i := 0; i < int(hPixels); i++ {
			data[i] = make([]float64, vPixels)
		}

		for h := 0; h < int(v.X); h++ {
			for x := 0; x < int(v.Y); x++ {
				du := dus[x*int(v.X)+h]
				for i := 0; i < 8; i++ {
					insertX := (i + int(h)*8) * int(xRange) //采样点少的情况下需要步偏移

					for j := 0; j < 8; j++ {
						insertY := (j + x*8) * int(yRange)
						value := du[i][j]
						for rx := 0; rx < int(xRange); rx++ {
							for ry := 0; ry < int(yRange); ry++ {
								data[insertX+rx][insertY+ry] = value
								if debug {
									fmt.Println("---ccc", value, insertX+rx, insertY)
								}
							}
						}
					}
				}
			}
		}

		index := v.Index
		for i := 0; i < int(hPixels); i++ {
			for j := 0; j < int(vPixels); j++ {
				_la := data[i][j] + float64(128)
				if debug {
					fmt.Println("---dddd", _la, data[i][j], i, j)
				}
				output[i][j][index] = _la
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
func ycrcb2rgb(YCrCb []float64) []int {
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

func keys(list []image_research.HtDataArr) []string {
	_list := make([]string, 0)
	for _, v := range list {
		_list = append(_list, v.Key)
	}
	return _list
}

func (a *JData) decodeHuffman(ctx context.Context, input []byte, colorComponentId uint, lastDc int) (int, []int) {
	cursor := 0
	output := make([]int, 0)
	for i := 0; i < 64; i++ {
		_r := a.colorInfo[colorComponentId].ACHTID
		var _d int
		if i == 0 {
			_r = a.colorInfo[colorComponentId].DCHTID
		} else {
			_d = 1
		}
		_key := fmt.Sprintf("%d-%d", _d, _r)
		_ht := a.Ht[_key].DataArr
		_htMap := a.Ht[_key].Data
		_keys := keys(_ht)
		//fmt.Println(i, "_key", _key, _d, _r)
		//fmt.Println("_ht", _ht)
		var value int
		for _, v := range _keys {
			length := len(v)
			//fmt.Println("v---", v, len(v))
			var subBuffer []byte
			if len(input) >= cursor+length {
				subBuffer = image_research.ReadBytesByStartAndEnd(input, uint(cursor), uint(cursor+length))
				//fmt.Println("subBuffer --", string(subBuffer), cursor, length)
			} else {
				subBuffer = image_research.SliceArr(input, cursor)
				tempSubBuffer := make([]byte, length-len(subBuffer))
				tempSubBuffer = image_research.FillBytes(tempSubBuffer, '1')
				subBuffer = image_research.Concat(subBuffer, tempSubBuffer, length)
				//fmt.Println("subBuffer 222---", string(subBuffer))
			}
			if string(subBuffer) == v {
				cursor += length
				value = _htMap[v].Value
				//fmt.Println("length", length, v, value, _keys)
				break
			}
		}
		//fmt.Println("bc---cur:", cursor)
		//fmt.Println("iiii--2", i)

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
				_bitData := image_research.ReadBytesByStartAndEnd(input, uint(cursor), uint(bitCount+cursor))
				_bitData1, _ := image_research.Str2Bin(string(_bitData))
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
			bitString := image_research.NumberToString(int64(value))
			zeroCount, _ := strconv.ParseInt(bitString[0:4], 2, 64) //数据前0的个数
			_bitCount, _ := strconv.ParseInt(bitString[4:8], 2, 64) //数据的位数
			//fmt.Println("bitCount: ", _bitCount, bitString)
			bitCount = int(_bitCount)
			//fmt.Println(value, string(bitString), "zeroCount", zeroCount, "_bitCount", _bitCount)
			if zeroCount == 0 && bitCount == 0 {
				//解析到 (0, 0) ，表示到达EOB，意味着后面的都是0
				for len(output) < 64 {
					output = append(output, 0)
				}
				//fmt.Println("break")
				break
			} else {
				if bitCount == 0 {
					bitData = 0
				} else {
					_bitData := image_research.ReadBytesByStartAndEnd(input, uint(cursor), uint(bitCount+cursor))
					_bitData1, _ := image_research.Str2Bin(string(_bitData))
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
		//fmt.Println("iiii--3", i)
		//fmt.Println("bc---",  bitCount, cursor, bitData)
		output = append(output, bitData)
		cursor += bitCount
		//fmt.Println("iiii--4", i)
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
	//遇到结束符号，结束
	if marker == 0xd9 {
		return nil
	}
	//读取2个字节
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
	a.height = image_research.ReadInt16(chunkData, 1)
	a.width = image_research.ReadInt16(chunkData, 3)
	a.colorComponentCount = chunkData[5]
	if a.colorComponentCount != 3 {
		return errors.New("仅支持YCrCb彩色图")
	}
	chunk := image_research.SliceArr(chunkData, 6)
	a.colorComponents = make(map[uint]ColorComponentsItem)
	for i := 0; i < int(a.colorComponentCount); i++ {
		colorComponentId := image_research.ReadInt8(chunk, 0)
		packField := image_research.ReadInt8(chunk, 1)
		_pF := image_research.NumberToString(int64(packField))
		_x, _ := strconv.ParseInt(_pF[0:4], 2, 10)
		_y, _ := strconv.ParseInt(_pF[4:8], 2, 10)
		qrId := image_research.ReadInt8(chunk, 2)
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
		chunk = image_research.SliceArr(chunk, 3)
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
		packedField := image_research.ReadInt8(chunkData, 0)
		_pf := image_research.NumberToString(int64(packedField))
		_type := _pf[3] //哈夫曼类型 0-DC直流 1-AC交流表
		_id := _pf[7]   //哈夫曼表id
		var length uint = 0
		countArr := make([]uint, 0)
		for i := 0; i < 16; i++ {
			_count := image_research.ReadInt8(chunkData, i+1)
			countArr = append(countArr, _count)
			length += _count
		}
		data, err := image_research.ReadBytes(chunkData, 17, length)
		if err != nil {
			return err
		}
		_data, _dataArr, _ := image_research.CreateHuffmanTree(data, countArr)
		_ht := Ht{
			Type:    image_research.AsciiToStr(int(_type)),
			Data:    _data,
			DataArr: _dataArr,
		}
		_key := fmt.Sprintf("%s-%s", image_research.AsciiToStr(int(_type)), image_research.AsciiToStr(int(_id)))
		a.Ht[_key] = _ht
		chunkData = image_research.SliceArr(chunkData, int(length)+17)
	}
	return nil
}

func (a *JData) decodeSOS(ctx context.Context, chunkData []byte) {
	a.colors = uint8(image_research.ReadInt8(chunkData, 0))
	chunkData = image_research.SliceArr(chunkData, 1)
	for i := 0; i < int(a.colors); i++ {
		colorId := image_research.ReadInt8(chunkData, 0) //1 - Y，2 - Cb，3 - Cr，4 - I，5 - Q
		packedField := image_research.ReadInt8(chunkData, 1)
		_pf := image_research.NumberToString(int64(packedField))

		_dchtId, _ := strconv.ParseInt(_pf[0:4], 2, 10) //DC哈夫曼表id
		_achtId, _ := strconv.ParseInt(_pf[4:8], 2, 10) //AC哈夫曼表id

		a.colorInfo[colorId] = ColorInfoItem{
			DCHTID: _dchtId,
			ACHTID: _achtId,
		}
		chunkData = image_research.SliceArr(chunkData, 2)
	}
	a.collectData = true //准备开始收集数据
}

//解码DQT，Define Quantization Table，定义量化表
func (a *JData) decodeDQT(ctx context.Context, chunkData []byte) error {
	for len(chunkData) > 0 {
		packedField := image_research.ReadInt8(chunkData, 0)
		_pf := image_research.NumberToString(int64(packedField))

		accuracy, _ := strconv.ParseInt(_pf[0:4], 2, 10) //量化表精度，0 - 1字节，1 - 2字节
		id, _ := strconv.ParseInt(_pf[4:8], 2, 10)       //量化表id，取值范围为0 - 3，所以最多可有4个量化表

		length := 64 * (accuracy + 1)
		data, err := image_research.ReadBytes(chunkData, 1, uint(length))
		if err != nil {
			return err
		}
		a.qt[uint(id)] = QtItem{
			accuracy: accuracy,
			data:     data,
		}

		chunkData = image_research.SliceArr(chunkData, len(chunkData)+1)
	}
	return nil
}

//解码DRI，Define Restart Interval，定义差分编码累计复位的间隔
func (a *JData) decodeDRI(ctx context.Context, chunkData []byte) {
	a.restartInterval = image_research.ReadInt16(chunkData, 0)
}

//Application，应用程序保留标记0
func (a *JData) decodeAPP0(ctx context.Context, chunkData []byte) error {
	_formatB, err := image_research.ReadBytes(chunkData, 0, 5)
	if err != nil {
		return err
	}
	a.format = image_research.BufferToString(_formatB)
	a.mainVersion = image_research.ReadInt8(chunkData, 5)
	a.subVersion = image_research.ReadInt8(chunkData, 6)

	a.unit = image_research.ReadInt8(chunkData, 7)
	a.xPixel = image_research.ReadInt16(chunkData, 8)
	a.yPixel = image_research.ReadInt16(chunkData, 10)

	a.thumbnailWidth = image_research.ReadInt8(chunkData, 12)
	a.thumbnailHeight = image_research.ReadInt8(chunkData, 13)

	if a.thumbnailWidth > 0 && a.thumbnailHeight > 0 && len(chunkData) > 14 {
		thumbnail := image_research.SliceArr(chunkData, 14)
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
