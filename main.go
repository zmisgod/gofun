package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	httpServer()
}

type RangeFileInfo struct {
	Size int64
	StartId int
	EndId int
	Length int
	IsApart bool
	File *os.File
}

func handleRangeFile(fileName string, rangeStr string) (*RangeFileInfo, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	if rangeStr == "" {
		fileObj := RangeFileInfo{
			Size:    size,
			StartId: 0,
			EndId:   int(size),
			Length:  int(size),
			IsApart: false,
			File:file,
		}
		return &fileObj, nil
	}
	startId, endId := parseRange(rangeStr)
	_, _ = file.Seek(int64(startId), io.SeekStart)
	if endId > int(size) {
		endId = int(size)
	}
	length := endId - startId+1
	fileObj := RangeFileInfo{
		Size:    size,
		StartId: startId,
		EndId:   endId,
		Length:  length,
		IsApart: true,
		File:file,
	}
	return &fileObj, nil
}

func httpServer() {
	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		fileName := "./image.jpg"
		getRange := request.Header.Get("Range")
		obj, err := handleRangeFile(fileName, getRange)
		if err != nil {
			SendError(writer, http.StatusForbidden, err.Error())
			return
		}
		defer obj.File.Close()
		if obj.IsApart {
			fmt.Printf("file size: %d startId %d - endId %d length %d\n", obj.Size, obj.StartId, obj.EndId, obj.Length)
			writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", fileName))
			writer.Header().Set("Content-Range", getContentRange(obj.StartId, obj.EndId, int(obj.Size)))
			writer.WriteHeader(http.StatusPartialContent)
			result := make([]byte, obj.Length)
			d, _ := obj.File.Read(result)
			fmt.Println(d)
			_, err = fmt.Fprint(writer, result)
			if err != nil {
				SendError(writer, http.StatusForbidden, err.Error())
				return
			}
		}else{
			writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", fileName))
			writer.WriteHeader(http.StatusOK)
			_, err = io.Copy(writer, obj.File)
			if err != nil {
				SendError(writer, http.StatusForbidden, err.Error())
				return
			}
		}
	})
	host:= "127.0.0.1:8212"
	fmt.Println("start http server at "+ host)
	if err := http.ListenAndServe(host, nil); err != nil {
		fmt.Println(err)
	}
}

func getImage(startId, endId int, outFile *os.File) error {
	fileName :=  "./image.jpg"
	obj, err := handleRangeFile(fileName, fmt.Sprintf("bytes=%d-%d", startId, endId))
	if err != nil {
		return err
	}
	defer obj.File.Close()
	result := make([]byte, endId-startId+1)
	_n, err := obj.File.Read(result)
	fmt.Printf("-- %d read %d\n", startId, _n)
	_, err = outFile.WriteAt(result, int64(startId))
	return err
}

func test() {
	dR := 10
	size := 3160785
	per :=  size/dR
	saveFile, _ := os.Create("res.jpg")
	saveFile.Truncate(int64(size))
	defer saveFile.Close()
	for i := 0; i < dR; i++ {
		startId := i * per
		endId := (i+1)*per - 1
		if i == (dR - 1) {
			endId = size
		}
		go func() {
			fmt.Println(startId, endId)
			if err := getImage(startId, endId, saveFile); err != nil {
				fmt.Println(err)
			}
		}()
	}
	time.Sleep(10 * time.Second)
}

func parseRange(rangeStr string) (int, int) {
	startId := 0
	endId := 0
	if rangeStr != "" {
		var re = regexp.MustCompile(`(?m)([0-9]*-[0-9]*)`)
		list := re.FindAllStringSubmatch(rangeStr, 50)
		if len(list) > 0 && len(list[0]) > 1 {
			split := strings.Split(list[0][1], "-")
			if len(split) > 1 {
				startId, _ = strconv.Atoi(split[0])
				endId, _ = strconv.Atoi(split[1])
			}
		}
	}
	return startId, endId
}

func getContentRange(startId, endId, size int) string {
	return fmt.Sprintf("bytes %d-%d/%d", startId, endId, size)
}

func SendError(writer http.ResponseWriter, code int, msg string) {
	fmt.Println(msg)
	writer.WriteHeader(code)
	_, _ = writer.Write([]byte(msg))
}
