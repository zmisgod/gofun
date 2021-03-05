package console_color

import "fmt"

type TextColor string

type BgColor string

const (
	TextBlack  TextColor = "30"
	TextRed    TextColor = "31"
	TextGreen  TextColor = "32"
	TextYellow TextColor = "33"
	TextBlue   TextColor = "34"
	TextPurple TextColor = "35"
	TextSky    TextColor = "36"
	TextWhite  TextColor = "37"
)

const (
	BgBlack  BgColor = "40"
	BgRed    BgColor = "41"
	BgGreen  BgColor = "42"
	BgYellow BgColor = "43"
	BgBlue   BgColor = "44"
	BgPurple BgColor = "45"
	BgSky    BgColor = "46"
	BgWhite  BgColor = "47"
)

type options struct {
	BgColor   string
	TextColor string
}

type Option func(*options)

func SetBgColor(color BgColor) Option {
	return func(o *options) {
		o.BgColor = string(color)
	}
}

func SetTextColor(color TextColor) Option {
	return func(o *options) {
		o.TextColor = string(color)
	}
}

func ConsoleColor(text string, colors ...Option) string {
	var con options
	for _, c := range colors {
		c(&con)
	}
	if con.TextColor == "" {
		con.TextColor = string(TextGreen)
	}
	colorString := ""
	if con.BgColor == "" {
		colorString = fmt.Sprintf("%sm", con.TextColor)
	} else {
		colorString = fmt.Sprintf("%s;%sm", con.BgColor, con.TextColor)
	}
	return fmt.Sprintf("\033[%s%s\033[0m", colorString, text)
}
