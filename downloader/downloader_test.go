package downloader

import (
	"fmt"
	"log"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	text := "https://www.imax.com/download/file/fid/16840"
	down, err := NewDownloader(text)
	if err != nil {
		log.Println(err)
	}else{
		err = down.SaveFile()
		if err != nil {
			fmt.Println(err)
		}else{
			fmt.Println(down.SaveName)
		}
	}
}
