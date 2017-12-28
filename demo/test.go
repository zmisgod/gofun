package main

import (
	"fmt"
	"strings"
)

func main() {
	url := "https://zmisgod.com/star.png"
	res := strings.Split(url, ".")
	fmt.Println(res[len(res)-1])
}
