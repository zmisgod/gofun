package range_server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ReadRangeFile interface {
	Read(ctx context.Context, startId, endId int) error
	Prepare(ctx context.Context) error
	Run(ctx context.Context) error
}

type ReadRangeCommon struct {
	FileName     string
	SaveFileName string
	fileSize     int
	ChunkSize    int
	outFile      *os.File
}

func (a *ReadRangeCommon) SetOutPutFileName(fileName string) {
	a.SaveFileName = fileName
}

func (a *ReadRangeCommon) SetChunkSize(size int) {
	a.ChunkSize = size
}

func (a *ReadRangeCommon) Close() error {
	return a.outFile.Close()
}

type RangeFileInfo struct {
	Size    int64
	StartId int
	EndId   int
	Length  int
	IsApart bool
	File    *os.File
}

func SimulateRangeServer(host string, port int, fileName string, isPart bool) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		getRange := request.Header.Get("Range")
		obj, err := handleRangeFile(fileName, getRange)
		if err != nil {
			SendError(writer, http.StatusForbidden, err.Error())
			return
		}
		defer obj.File.Close()
		if isPart {
			writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", fileName))
			writer.Header().Set("Content-Range", getContentRange(obj.StartId, obj.EndId, int(obj.Size)))
			writer.WriteHeader(http.StatusPartialContent)
			result := make([]byte, obj.Length)
			_, err := obj.File.Read(result)
			if err != nil {
				fmt.Println(err)
			}
			_, _ = writer.Write(result)
			if err != nil {
				SendError(writer, http.StatusForbidden, err.Error())
				return
			}
		} else {
			writer.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", fileName))
			writer.Header().Set("Content-Length", fmt.Sprintf("%d", obj.Size))
			writer.WriteHeader(http.StatusOK)
			_, err = io.Copy(writer, obj.File)
			if err != nil {
				SendError(writer, http.StatusForbidden, err.Error())
				return
			}
		}
	})
	serverHost := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("start http server at " + serverHost)
	if err := http.ListenAndServe(serverHost, nil); err != nil {
		fmt.Println(err)
	}
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
			File:    file,
		}
		return &fileObj, nil
	}
	startId, endId := parseRange(rangeStr)
	_, _ = file.Seek(int64(startId), io.SeekStart)
	if endId > int(size) {
		endId = int(size)
	}
	length := endId - startId + 1
	fileObj := RangeFileInfo{
		Size:    size,
		StartId: startId,
		EndId:   endId,
		Length:  length,
		IsApart: true,
		File:    file,
	}
	return &fileObj, nil
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
