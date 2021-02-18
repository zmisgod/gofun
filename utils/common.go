package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CreateFile(textName string) {
	_, err := os.Stat(textName)
	if err != nil {
		file, err := os.Create(textName)
		CheckError(err)
		defer file.Close()
	}
}

//CreateFileReError 创建文件
func CreateFileReError(textName string) (*os.File, error) {
	_, err := os.Stat(textName)
	if err != nil {
		file, err := os.Create(textName)
		if err != nil {
			return nil, err
		}
		return file, nil
	}else{
		return nil, errors.New(fmt.Sprintf("file %s exits", textName))
	}
}

//CreateFolder 创建文件夹
func CreateFolder(folderName string) (bool, error) {
	checkFolderNotExists := CheckPathIsNotExists(folderName)
	if checkFolderNotExists {
		err := os.MkdirAll(folderName, 0777)
		if err != nil {
			return false, err
		}
		log.Printf("create floder %s successful\n", folderName)
		return true, nil
	}
	return false, nil
}

//CheckPathIsNotExists 检查文件是否存在 返回true 不存在， false 存在
func CheckPathIsNotExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}