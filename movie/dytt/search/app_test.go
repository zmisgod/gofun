package search

import (
	"context"
	"testing"
)

func TestFetch(t *testing.T) {
	obj, err := Fetch(context.Background(), "旺达幻视")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj)
}
