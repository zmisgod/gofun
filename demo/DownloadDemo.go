package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	url := []string{"http://www.logosc.cn/Public/imgs/touxiang.jpg"}
	floderName := "./images/" + "demo"
	err := createFolder(floderName)
	if err != nil {
		panic(err)
	}
	count := make(chan int)
	go downloadImage(url, floderName, "demo1", "jpg", count)
	finish := <-count
	fmt.Printf("finish download %d image\n", finish)
}

/**
 * 创建文件夹
 */
func createFolder(floderName string) error {
	checkFloderNotExists, err := checkPathIsNotExists(floderName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if checkFloderNotExists {
		err := os.MkdirAll(floderName, 0777)
		if err != nil {
			return err
		}
		fmt.Printf("create floder %s successful\n", floderName)
	} else {
		fmt.Printf("floder %s already exists\n", floderName)
	}
	return nil
}

/**
 * 检查文件是否存在
 * 返回true 不存在， false 存在
 */
func checkPathIsNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 下载图片
 */
func downloadImage(urls []string, folderPath string, fileName string, fileType string, count chan int) {
	ct := 0
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		imgByte, _ := ioutil.ReadAll(resp.Body)
		notExist, _ := checkPathIsNotExists(folderPath + "/" + fileName + "." + fileType)
		if notExist {
			fp, _ := os.Create(folderPath + "/" + fileName + "." + fileType)
			defer fp.Close()
			fp.Write(imgByte)
			ct++
			fmt.Printf("create %s/%s.%s successful\n", folderPath, fileName, fileType)
		} else {
			fmt.Printf("sorry %s/%s.%s exist\n", folderPath, fileName, fileType)
		}
	}
	count <- ct
}
