package fileserver

import (
	"fmt"
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	port := 12345
	host := ""
	showTem := "tem"
	template := ""
	basePath := "./"
	validDownload := []string{"svg", "jpg", "png", "jpeg", "gif"}
	dir, err := Create(port, basePath, host, showTem, template, validDownload)
	if err != nil {
		log.Fatal(err)
	}
	if err := dir.CreateServer(); err != nil {
		fmt.Println(err)
	}
}
