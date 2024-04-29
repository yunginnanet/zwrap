package zwrap

import (
	"fmt"
	"github.com/rs/zerolog"
	"strings"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

var LevelColors = map[zerolog.Level]int{
	zerolog.TraceLevel: colorMagenta,
	zerolog.DebugLevel: colorYellow,
	zerolog.InfoLevel:  colorGreen,
	zerolog.WarnLevel:  colorRed,
	zerolog.ErrorLevel: colorRed,
	zerolog.FatalLevel: colorRed,
	zerolog.PanicLevel: colorRed,
}

// Colorize returns the string s wrapped in ANSI code c, unless disabled is true or c is 0.
func Colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// FormattedLevels are used by ConsoleWriter's consoleDefaultFormatLevel
// for a short level name.
var FormattedLevels = map[zerolog.Level]string{
	zerolog.TraceLevel: "TRC",
	zerolog.DebugLevel: "DBG",
	zerolog.InfoLevel:  "INF",
	zerolog.WarnLevel:  "WRN",
	zerolog.ErrorLevel: "ERR",
	zerolog.FatalLevel: "FTL",
	zerolog.PanicLevel: "PNC",
}

func LogLevelFmt(noColor bool) zerolog.Formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			level, _ := zerolog.ParseLevel(ll)
			fl, ok := FormattedLevels[level]
			if ok {
				l = Colorize(fl, LevelColors[level], noColor)
			} else {
				l = strings.ToUpper(ll)[0:3]
			}
		} else {
			if i == nil {
				l = "???"
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}
