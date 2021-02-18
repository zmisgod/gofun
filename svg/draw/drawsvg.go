package draw

import (
	"errors"
	"fmt"
	"strconv"
)

//Canvas 画布信息
type Canvas struct {
	viewBox    viewBox
	reSize     int
	longResize int
	laResize   int
}

//ViewBox ViewBox
type viewBox struct {
	widthOne  int
	heightOne int
	widthTwo  int
	heightTwo int
}

//Circle 画圆点
type circle struct {
	radius      int
	stroke      string
	fill        string
	strokeWidth int
}

//path指令，如果不传则默认直线
var directive = []int{
	0, //直线
	1, //曲线
}

//IPath IPath
type IPath struct {
	Group int //0 1 2
	Long  string
	Lat   string
	ID    int
	Name  string

	//下面仅对PATH有效
	Directive int //指令
}

//Path 路径
type Path struct {
	MaxGroup    int //最大的Group是多少
	PathInfo    []IPath
	d           []string
	g           string
	Fill        string
	Stroke      string
	StrokeWidth int    `json:"stroke-width"`
	FillOpacity string `json:"fill-opacity"`
	Aid         string //额外信息1
	Alt         string //额外信息2
}

//Output 输出
type Output struct {
	Name   string
	Floder string
}

//SVG svg数据结构
type SVG struct {
	Canvas    Canvas `json:"canvas"` //svg画布与基础参数
	circle    circle //画点
	Output    Output `json:"output"` //输出
	pathLists []Path
}

//Create 创建一个SVG
func Create() *SVG {
	var svg SVG
	svg.Canvas.longResize = 18
	svg.Canvas.laResize = 70
	return &svg
}

//SetResize 设置resize
func (svg *SVG) SetResize(resize int) {
	svg.Canvas.reSize = resize
}

//SetCircle 设置circle形状
func (svg *SVG) SetCircle(radius, strokeWidth int, stroke, fill string) {
	if radius <= 0 {
		svg.circle.radius = 1
	} else {
		svg.circle.radius = radius
	}
	if stroke == "" {
		svg.circle.stroke = "blue"
	} else {
		svg.circle.stroke = stroke
	}
	if fill == "" {
		svg.circle.fill = "transparent"
	} else {
		svg.circle.fill = fill
	}
	if strokeWidth <= 0 {
		svg.circle.strokeWidth = 1
	} else {
		svg.circle.strokeWidth = strokeWidth
	}
}

//SetOutputFileName 设置output文件名称
func (svg *SVG) SetOutputFileName(floder, name string) {
	svg.Output.Name = name
	svg.Output.Floder = floder
}

//beforeDraw 画图前检查
func (svg *SVG) beforeDraw() error {
	if svg.Output.Floder == "" {
		return errors.New("please set the output floder first")
	}
	if svg.Output.Name == "" {
		svg.Output.Name = "zmis.me.svg"
	}
	return nil
}

//SetPath 设置path
func (svg *SVG) SetPath(path Path) {
	svg.pathLists = append(svg.pathLists, path)
}

//Draw 画操作
func (svg *SVG) Draw() string {
	for index, path := range svg.pathLists {
		content := ""
		//画线段
		path.DrawSVG(svg)

		content = path.String()
		if len(path.PathInfo) > 0 {
			for _, v := range path.PathInfo {
				longStr, latiStr := svg.dataConversion(v.Lat, v.Long)
				content += svg.circle.String(latiStr, longStr, v.Name, v.ID)
			}
		}
		svg.pathLists[index].g = "<g>" + content + "</g>"
	}

	svg.Canvas.viewBox.widthTwo += 100
	svg.Canvas.viewBox.heightTwo += 100
	return svg.String()
}

//DrawSVG 画d
func (path *Path) DrawSVG(svg *SVG) {
	pathInfoCount := len(path.PathInfo)
	if pathInfoCount > 0 {
		pathStr := ""
		for k, v := range path.PathInfo {
			prefix := "L "
			if k == 0 {
				prefix = "M"
			}
			if v.Directive == 1 {
				prefix = "Q "
			}
			pathStr += prefix

			if v.Directive == 1 && k >= 1 {
				centerPointLa := path.PathInfo[k-1].Long
				centerPointLo := v.Lat
				longStr, latiStr := svg.dataConversion(centerPointLa, centerPointLo)
				pathStr += longStr + " " + latiStr + " "
			}

			longStr, latiStr := svg.dataConversion(v.Lat, v.Long)
			pathStr += longStr + " " + latiStr + " "
		}
		path.d = append(path.d, pathStr)
	}
}

//dataConversion 数据转换
func (svg *SVG) dataConversion(data1, data2 string) (string, string) {
	f1, _ := strconv.ParseFloat(data1, 32)
	res1 := int((f1 - float64(svg.Canvas.laResize)) * float64(svg.Canvas.reSize))

	f2, _ := strconv.ParseFloat(data2, 32)
	res2 := int((f2 - float64(svg.Canvas.longResize)) * float64(svg.Canvas.reSize))

	if res1 > svg.Canvas.viewBox.widthTwo {
		svg.Canvas.viewBox.widthTwo = res1
	}
	if res2 > svg.Canvas.viewBox.heightTwo {
		svg.Canvas.viewBox.heightTwo = res2
	}
	str1 := strconv.Itoa(res1)
	str2 := strconv.Itoa(res2)
	return str1, str2
}

//String path
func (path *Path) String() string {
	content := ""
	if len(path.d) > 0 {
		for _, d := range path.d {
			content += fmt.Sprintf("<path d=\"%s\" stroke=\"%s\" fill=\"%s\" stroke-width=\"%d\" fill-opacity=\"%s\" alt=\"%s\" aid=\"%s\" />", d, path.Stroke, path.Fill, path.StrokeWidth, path.FillOpacity, path.Alt, path.Aid)
		}
	}
	return content
}

//String svg
func (svg *SVG) String() string {
	content := ""
	if len(svg.pathLists) > 0 {
		for _, path := range svg.pathLists {
			content += path.g
		}
	}
	return fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"%s\" version=\"1.1\"  width=\"100%%\" height=\"100%%\">%s</svg>", svg.Canvas.viewBox.String(), content)
}

//String circle
func (c *circle) String(long, la, name string, id int) string {
	return fmt.Sprintf("<circle cx=\"%s\" cy=\"%s\" r=\"%d\" stroke=\"%s\" fill=\"%s\" alt=\"%s\" aid=\"%d\" />", la, long, c.radius, c.stroke, c.fill, name, id)
}

//String viewBox
func (c *viewBox) String() string {
	w1 := strconv.Itoa(c.widthOne)
	h1 := strconv.Itoa(c.heightOne)
	w2 := strconv.Itoa(c.widthTwo)
	h2 := strconv.Itoa(c.heightTwo)
	return w1 + "," + h1 + "," + w2 + "," + h2
}
