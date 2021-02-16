package bloom_filter

import (
	"testing"
)

func TestNewBloomFilter(t *testing.T) {
	bf := NewBloomFilter()
	bf.Add("1")
	bf.Add("2")
	bf.Add("3")
	if !bf.Contains("1") {
		t.Error("1")
	}
	if !bf.Contains("3") {
		t.Error("3")
	}
	if bf.Contains("13") {
		t.Error("13")
	}
}