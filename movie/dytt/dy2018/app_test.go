package dy2018

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestFetchByID(t *testing.T) {
	obj, err := FetchByID(context.Background(), "103070")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(obj)
}

func TestSearch(t *testing.T) {
	res, err := SearchMovies(context.Background(), "告白")
	if err != nil {
		log.Fatalln(err)
	}
	if res != nil {
		for _, v := range res.List {
			fmt.Println(v)
		}
		fmt.Println(res.Page)
	}
}
