package fetchonline

import (
	"fmt"
	"os"
)

//FetchOnline 从远程下载
type FetchOnline struct {
	TargetURL      string
	DownloadFolder string
}

//SetDownloadFolder 设置下载目录
func (t *FetchOnline) SetDownloadFolder(path string) {

}

//StartServer 开启下载服务器
func (t *FetchOnline) StartServer() {

}

//CheckFolder 检查文件夹是否存在
func (t *FetchOnline) CheckFolder() (bool, error) {
	checkFolderNotExists, err := checkPathIsNotExists(t.DownloadFolder)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if checkFolderNotExists {
		err := os.MkdirAll(t.DownloadFolder, 0777)
		if err != nil {
			return false, err
		}
		fmt.Printf("create floder %s successful\n", t.DownloadFolder)
		return true, nil
	}
	return false, err
}

//Download 下载
func (t *FetchOnline) Download() {

}

//CheckQueue 检查redis队列
func (t *FetchOnline) CheckQueue() {

}

//HandleQueue 处理队列数据
func (t *FetchOnline) HandleQueue() {

}
