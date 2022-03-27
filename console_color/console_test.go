package console_color

import (
	"fmt"
	"testing"
)

func TestConsoleColor(t *testing.T) {
	text := ConsoleColor("22222",
		SetBgColor(BgYellow),
		SetTextColor(TextWhite),
	)
	fmt.Println(text)
}
