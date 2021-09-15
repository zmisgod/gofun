package tencent

import (
	"os"
	"testing"
)

func TestNewJson(t *testing.T) {
	_file, err := os.Open("./res.json")
	if err != nil {
		t.Error(err)
	}
	re, err := NewJson(_file)
	if err != nil {
		t.Error(err)
	}
	re.handleData()
}