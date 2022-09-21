package dy2018

import (
	"context"
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	ctx := context.Background()
	res, _ := SearchMovies(ctx, "断桥")
	if res != nil {
		for res.HasMore(ctx) {
			list, _ := res.Next(ctx)
			for _, v := range list {
				fmt.Println(v)
				dLinks, _ := v.GetDownloadUrls(ctx)
				fmt.Println(dLinks)
				fmt.Println("----------")
			}
		}
	}
}
