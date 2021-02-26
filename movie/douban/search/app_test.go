package search

import (
	"context"
	"testing"
)

func TestFindMovieInfo(t *testing.T) {
	obj, err := Fetch(context.Background(), "34841067")
	if err != nil {
		t.Fatal(err)
	}else{
		t.Log(obj)
	}
}