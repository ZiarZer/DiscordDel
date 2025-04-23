package utils

import (
	"fmt"
	"strings"
)

type LogLevel struct {
	Name          string
	BashColorCode string
}

var (
	DEBUG   = &LogLevel{Name: "DEBUG", BashColorCode: "1"}
	INFO    = &LogLevel{Name: "INFO", BashColorCode: "1;44"}
	SUCCESS = &LogLevel{Name: "SUCCESS", BashColorCode: "1;42"}
	WARNING = &LogLevel{Name: "WARNING", BashColorCode: "1;31;103"}
	ERROR   = &LogLevel{Name: "ERROR", BashColorCode: "1;41"}
	FATAL   = &LogLevel{Name: "FATAL", BashColorCode: "1;100"}
)

func cJust(s string, length int, fill string) string {
	leftLen := (length - len(s)) / 2
	rightLen := length - len(s) - leftLen
	return strings.Repeat(fill, leftLen) + s + strings.Repeat(fill, rightLen)
}

func colorText(text string, bashColorCode string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", bashColorCode, text)
}

func getLogTag(logLevel *LogLevel) string {
	if logLevel == nil {
		return "         "
	}
	return colorText(fmt.Sprintf("[%s]", colorText(cJust(logLevel.Name, 7, " "), logLevel.BashColorCode)), "1")
}

func InternalLog(message string, logLevel *LogLevel) {
	fmt.Printf("%s %s\n", getLogTag(logLevel), message)
}
