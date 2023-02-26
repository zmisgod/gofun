package png

import (
	"context"
	"fmt"
	"testing"
)

func TestNewPng(t *testing.T) {
	obj, err := NewPng("./test.png")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(obj.decode(context.Background()))
}