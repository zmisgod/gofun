package imdb

import (
	"context"
	"testing"
)

func TestFetch(t *testing.T) {
	obj, err := Fetch(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj)
}
