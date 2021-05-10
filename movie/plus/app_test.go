package plus

import (
	"fmt"
	"strings"
	"testing"
)

func TestAsdse(t *testing.T) {
	fmt.Println( strings.Split(strings.Trim(",1,2,34,5,,,,,,", ","), ","))
}