package downloader

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	text := "https://www.imax.com/download/file/fid/16840"
	down, err := NewDownloader(text, SetTimeout(60), SetBreakPointContinueUpload(true))
	if err != nil {
		log.Println(err)
	} else {
		err = down.SaveFile(context.Background())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(down.SaveName)
		}
	}
}
