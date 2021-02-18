package fileserver

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	port := 12345
	host := ""
	showTem := "tem"
	template := ""
	basePath := "./../pictureSpider/"
	validDownload := []string{"svg", "jpg", "png", "jpeg", "gif"}
	dir := Create(port, basePath, host, showTem, template, validDownload)
	if err := dir.CreateServer(); err != nil {
		fmt.Println(err)
	}
}