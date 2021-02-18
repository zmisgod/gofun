package utils

import (
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

/**
 * 创建文件夹
 */
func CreateFolder(folderName string) (bool, error) {
	checkFolderNotExists, err := CheckPathIsNotExists(folderName)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if checkFolderNotExists {
		err := os.MkdirAll(folderName, 0777)
		if err != nil {
			return false, err
		}
		log.Printf("create floder %s successful\n", folderName)
		return true, nil
	}
	return false, err
}

/**
 * 检查文件是否存在
 * 返回true 不存在， false 存在
 */
func CheckPathIsNotExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}