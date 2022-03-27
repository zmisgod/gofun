package range_server

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type ReadLocalRangeFile struct {
	ReadRangeCommon
}

func (a *ReadLocalRangeFile) Read(ctx context.Context, startId, endId int) error {
	obj, err := handleRangeFile(a.FileName, fmt.Sprintf("bytes=%d-%d", startId, endId))
	if err != nil {
		return err
	}
	defer obj.File.Close()
	result := make([]byte, endId-startId+1)
	_n, err := obj.File.Read(result)
	fmt.Printf("-- %d read %d\n", startId, _n)
	_, err = a.outFile.WriteAt(result, int64(startId))
	return err
}

func (a *ReadLocalRangeFile) Prepare(ctx context.Context) error {
	if a.SaveFileName == "" {
		exp := strings.Split(a.FileName, ".")
		if len(exp) <= 1 {
			a.SaveFileName = fmt.Sprintf("%s%d", exp[len(exp)], time.Now().Unix())
		} else {
			exp[len(exp)-2] = fmt.Sprintf("%s%d", exp[len(exp)-2], time.Now().Unix())
			a.SaveFileName = strings.Join(exp, ".")
		}
	}
	if a.ChunkSize == 0 {
		a.ChunkSize = 10
	}
	fileOpen, err := os.Open(a.FileName)
	if err != nil {
		return err
	}
	_info, err := fileOpen.Stat()
	if err != nil {
		return err
	}
	a.fileSize = int(_info.Size())
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

func (a *ReadLocalRangeFile) Run(ctx context.Context) error {
	if err := a.Prepare(ctx); err != nil {
		return err
	}
	defer a.outFile.Close()
	per := a.fileSize / a.ChunkSize
	var wg sync.WaitGroup
	for i := 0; i < a.ChunkSize; i++ {
		startId := i * per
		endId := (i+1)*per - 1
		if i == (a.ChunkSize - 1) {
			endId = a.fileSize
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(startId, endId)
			if err := a.Read(ctx, startId, endId); err != nil {
				fmt.Println(err)
			}
		}()
	}
	wg.Wait()
	return nil
}
