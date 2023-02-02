package main

import (
	"github.com/zmisgod/gofun/image_research/jpeg"
	"log"
)

func main() {
	_, err := jpeg.NewFile("./image_research/jpeg/test.jpeg")
	if err != nil {
		log.Fatalln(err)
	}
}