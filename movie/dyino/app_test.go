package dyino

import (
	"context"
	"log"
	"testing"
)

func TestNewMovieList(t *testing.T) {
	obj, err := Fetch(context.Background(), 2020, PublishVersionIMaxStereo, 100)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(obj)
}
