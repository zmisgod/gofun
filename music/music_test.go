package music

import (
	"context"
	"fmt"
	"testing"
)

func TestNewFetchMusic(t *testing.T) {
	res, err := NewFetchMusic("王心凌《大眠》 https://c.y.qq.com/base/fcgi-bin/u?__=d3IYVRj @QQ音乐")
	if err != nil {

	}else{
		resUrl, err := res.GetDownloadURL(context.Background())
		fmt.Println(resUrl, err)
	}
}
