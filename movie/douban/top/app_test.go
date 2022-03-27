package top

import (
	"context"
	"testing"
)

func TestFetch(t *testing.T) {
	obj, err := Fetch(context.Background(), DouBanCategoryAction, 0, 20)
	if err != nil {
		t.Log(err)
	}
	t.Log(obj)
}
