package main

import (
	"fmt"

	"github.com/zmisgod/goTool/drawsvg"
)

func main() {
	svg := drawsvg.Create()
	svg.SetCircle(2, 1, "#42526e", "transparent")

	var path drawsvg.Path
	path.Fill = "transparent"
	path.Stroke = "blue"
	path.StrokeWidth = 2
	path.FillOpacity = "0.4"
	path.Aid = "121"
	path.Alt = "121212"
	path.MaxGroup = 1
	for i := 0; i < 3; i++ {
		var ipath drawsvg.IPath
		if i == 0 {
			ipath.Lat = "100"
			ipath.Long = "105"
			ipath.ID = 11111
			ipath.Name = "station1"
		} else if i == 1 {
			ipath.Lat = "120"
			ipath.Long = "125"
			ipath.Directive = 1
			ipath.ID = 222222
			ipath.Name = "station2"
		} else if i == 2 {
			ipath.Lat = "150"
			ipath.Long = "90"
			ipath.ID = 33333
			ipath.Name = "station3"
		}
		path.PathInfo = append(path.PathInfo, ipath)
	}
	svg.SetPath(path)
	content := svg.Draw()
	fmt.Println(content)
}
