package name

import (
	"context"
	"fmt"
	"testing"
)

func TestSearchMovieInfo(t *testing.T) {
	ctx := context.Background()
	obj, err := Fetch(ctx, "勇敢者游戏2：再战巅峰")
	if err != nil {
		t.Fatal(err)
	} else {
		for _, v := range obj {
			fmt.Println(v)
		}
	}
}
