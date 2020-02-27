package main

import fileserver "github.com/zmisgod/goTool/fileServer"

func main() {
	port := 12345
	host := ""
	showTem := "tem"
	template := ""
	basePath := "./../pictureSpider/"
	validDownload := []string{"svg", "jpg", "png", "jpeg", "gif"}
	dir := fileserver.Create(port, basePath, host, showTem, template, validDownload)
	dir.CreateServer()
}
