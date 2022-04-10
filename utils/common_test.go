package utils

import (
	"fmt"
	"testing"
)

func Test_RemoveMoneySuffix(t *testing.T) {
	fmt.Println(RemoveMoneySuffix("12.0"))
	fmt.Println(RemoveMoneySuffix("12.00"))
	fmt.Println(RemoveMoneySuffix("12"))
}