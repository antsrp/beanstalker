package visual

import "github.com/TwiN/go-color"

var (
	ColorBold   = color.Bold
	ColorCyan   = color.Cyan
	ColorRed    = color.Red
	ColorGreen  = color.Green
	ColorYellow = color.Yellow
	ColorPurple = color.Purple
	ColorBlack  = color.Black
	ColorBlue   = color.Blue
	ColorGray   = color.Gray
	ColorWhite  = color.White
	ColorRedBG  = color.RedBackground
)

func Colorize(clr string, s string) string {
	return color.Ize(clr, s)
}
