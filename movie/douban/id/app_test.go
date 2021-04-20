package id

import (
	"context"
	"github.com/zmisgod/gofun/downloader"
	"github.com/zmisgod/gofun/utils"
	"testing"
	"time"
)

func TestFindMovieInfo(t *testing.T) {
	ctx := context.Background()
	obj, err := Fetch(ctx, "25824686")
	if err != nil {
		t.Fatal(err)
	}else{
		if len(obj.CoverUrls) > 0 {
			for _,v := range obj.CoverUrls {
				d, err := downloader.NewDownloader(v, downloader.SetTimeout(10))
				if err == nil {
					if err := d.SaveFile(ctx); err != nil {
						t.Log(err)
					}
				}
				time.Sleep(time.Duration(utils.Rand(2, 10)) * time.Second)
			}
		}
		t.Log(obj)
	}
}