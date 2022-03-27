package dy2018

import (
	"context"
	"testing"
)

func TestFetchByID(t *testing.T) {
	obj, err := FetchByID(context.Background(), "103070")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj)
}
