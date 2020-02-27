
```
func main (){
    svg, err := drawsvg.Create()
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }

    var path drawsvg.Path
    path.Fill = "transparent"
    path.Stroke = "blue"
    path.StrokeWidth = 5
    path.FillOpacity = "0.4"
    path.Aid = "121"
    path.Alt = "121212"
    for i := 0; i < 3; i++ {
        var ipath drawsvg.IPath
        if i == 0 {
            ipath.Lat = "100"
            ipath.Long = "105"
        } else if i == 1 {
            ipath.Lat = "120"
            ipath.Long = "125"
            ipath.Directive = "wheel"
            ipath.WheelLa = "130"
            ipath.WheelLong = "235"
        } else if i == 2 {
            ipath.Lat = "80"
            ipath.Long = "76"
        }
        path.PathInfo = append(path.PathInfo, ipath)
    }
    svg.Path = path
    content := svg.Draw()
    fmt.Println(content)
}
```