package range_server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/zmisgod/gofun/utils"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReadRemoteRangeFile struct {
	TargetUrl string
	ReadRangeCommon
}

func (a *ReadRemoteRangeFile) Read(ctx context.Context, startId, endId int) error {
	rangeStr := utils.GetHeaderRange(startId, endId)
	resp, err := httpClient(ctx, a.TargetUrl, rangeStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.New(fmt.Sprintf("startId %d- endId %d response status is not valid,now is %d", startId, endId, resp.StatusCode))
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = a.outFile.WriteAt(result, int64(startId))
	return err
}

func httpClient(ctx context.Context, targetURL string, rangeStr string) (*http.Response, error) {
	client := &http.Client{}
	transPort := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = transPort
	request, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	if rangeStr != "" {
		request.Header.Set(utils.HeaderRange, rangeStr)
	}
	request = request.WithContext(ctx)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ReadRemoteRangeFile) Run(ctx context.Context) error {
	if err := a.Prepare(ctx); err != nil {
		return err
	}
	defer a.outFile.Close()
	per := a.fileSize / a.ChunkSize
	//var wg sync.WaitGroup
	for i := 0; i < a.ChunkSize; i++ {
		startId := i * per
		endId := (i+1)*per - 1
		if i == (a.ChunkSize - 1) {
			endId = a.fileSize
		}
		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		fmt.Println("Run-----", startId, endId)
		if err := a.Read(ctx, startId, endId); err != nil {
			fmt.Println("ReadError", err)
		}
		//}()
	}
	//wg.Wait()
	return nil
}

func (a *ReadRemoteRangeFile) Prepare(ctx context.Context) error {
	if err := a.GetBaseInfo(ctx); err != nil {
		return err
	}
	needRename := false
	if a.SaveFileName == "" || a.SaveFileName == a.FileName {
		exp := strings.Split(a.FileName, ".")
		if len(exp) <= 1 {
			needRename = true
			a.SaveFileName = fmt.Sprintf("%s%d", exp[len(exp)], time.Now().Unix())
		} else {
			needRename = true
			exp[len(exp)-2] = fmt.Sprintf("%s%d", exp[len(exp)-2], time.Now().Unix())
			a.SaveFileName = strings.Join(exp, ".")
		}
	}
	if !needRename {
		exp := strings.Split(a.SaveFileName, ".")
		if len(exp) <= 1 {
			needRename = true
			a.SaveFileName = fmt.Sprintf("%s%d", exp[len(exp)], time.Now().Unix())
		} else {
			needRename = true
			exp[len(exp)-2] = fmt.Sprintf("%s%d", exp[len(exp)-2], time.Now().Unix())
			a.SaveFileName = strings.Join(exp, ".")
		}
	}
	if a.ChunkSize == 0 {
		a.ChunkSize = 10
	}
	file, err := os.Create(a.SaveFileName)
	if err != nil {
		return err
	}
	a.outFile = file
	if err := a.outFile.Truncate(int64(a.fileSize)); err != nil {
		return err
	}
	return nil
}

func (a *ReadRemoteRangeFile) GetBaseInfo(ctx context.Context) error {
	resp, err := httpClient(ctx, a.TargetUrl, utils.GetHeaderRange(0, 1))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.New("status code is not 200 or 206")
	}
	if resp.StatusCode != http.StatusPartialContent {
		a.ChunkSize = 1
	}
	fileName := resp.Header.Get("Content-Disposition")
	if fileName != "" {
		var re = regexp.MustCompile(`(?m)filename="(.*)"`)
		rest := re.FindAllStringSubmatch(fileName, 1000)
		if len(rest) > 0 && len(rest[0]) > 0 {
			a.SaveFileName = rest[0][1]
		}
	}
	sizeInfo := resp.Header.Get("Content-Range")
	if sizeInfo != "" {
		var re = regexp.MustCompile(`(?m)bytes (.*)`)
		rest := re.FindAllStringSubmatch(sizeInfo, 1000)
		if len(rest) > 0 && len(rest[0]) > 0 {
			expSizeStr := strings.Split(rest[0][1], "/")
			if len(expSizeStr) >= 2 {
				a.fileSize, _ = strconv.Atoi(expSizeStr[1])
			}
		}
	}
	if a.fileSize == 0 {
		_size := resp.Header.Get("content-length")
		if _size != "" {
			a.fileSize, _ = strconv.Atoi(_size)
		}
		if a.fileSize == 0 {
			_size := resp.Header.Get("Content-Length")
			if _size != "" {
				a.fileSize, _ = strconv.Atoi(_size)
			}
		}
	}
	return nil
}
