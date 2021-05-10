package range_server

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

func TestReadLocalRangeFile(t *testing.T) {
	readRangeFile := ReadLocalRangeFile{
		ReadRangeCommon{
			FileName: "./image.jpg",
		},
	}
	if err := readRangeFile.Run(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestReadRemoteRangeFile(t *testing.T) {
	port := 12121
	go func() {
		SimulateRangeServer("127.0.0.1", port, "image.jpg", true)
	}()
	time.Sleep(1*time.Second)
	readRangeFile := ReadRemoteRangeFile{
		TargetUrl: fmt.Sprintf("http://127.0.0.1:%d/", port),
	}
	if err := readRangeFile.Run(ctx); err != nil {
		t.Fatal(err)
	}
}