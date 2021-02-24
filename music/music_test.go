package music

import (
	"context"
	"testing"
)

func TestNewFetchMusic(t *testing.T) {
	//netEastUrl := "分享五月天的单曲《小时候》https://y.music.163.com/m/song?id=187283&userid=271602530&app_version=8.1.30 (@网易云音乐)"
	qqMusicUrl := "m.o.v.e《Around The World》(《头文字D》 TV动画片头曲) https://c.y.qq.com/base/fcgi-bin/u?__=vro0fZR @QQ音乐"
	res, err := NewFetchMusic(qqMusicUrl)
	if err != nil {
		t.Fatal(err)
	}
	resUrl, err := res.GetDownloadURL(context.Background())
	if err != nil {
		t.Fatal(err)
	}else{
		t.Log(resUrl)
	}
}
