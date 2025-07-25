package color

type Color int

const (
	None Color = iota
	Black
	Blue
	Red
	Green
	Yellow
	Magenta
	Cyan
	White
)

var terminalColorCodes = map[Color]string{
	Black:   "\033[30m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	White:   "\033[37m",
}

const Reset = "\033[0m"

func (c Color) Paint(text string) string {
	colorCode, ok := terminalColorCodes[c]
	if !ok {
		colorCode = ""
	}

	return colorCode + text + Reset
}
