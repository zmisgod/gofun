package downloader

import (
	"fmt"
	"testing"
)

func TestNewDownloader(t *testing.T) {
	down, err := NewDownloader("https://cdn.poizon.com/leap/A5CEF94C-5BA4-45F4-BB9F-8B7E7ADA02C2.mov_dgTLUAnsBO.mp4")
	fmt.Println(err, down)
}
