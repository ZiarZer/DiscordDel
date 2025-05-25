package utils

import (
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
)

var (
	DEBUG   = &types.LogLevel{Name: "DEBUG", BashColorCode: "1"}
	INFO    = &types.LogLevel{Name: "INFO", BashColorCode: "1;44"}
	SUCCESS = &types.LogLevel{Name: "SUCCESS", BashColorCode: "1;42"}
	WARNING = &types.LogLevel{Name: "WARNING", BashColorCode: "1;31;103"}
	ERROR   = &types.LogLevel{Name: "ERROR", BashColorCode: "1;41"}
	FATAL   = &types.LogLevel{Name: "FATAL", BashColorCode: "1;100"}
)

func cJust(s string, length int, fill string) string {
	leftLen := (length - len(s)) / 2
	rightLen := length - len(s) - leftLen
	return strings.Repeat(fill, leftLen) + s + strings.Repeat(fill, rightLen)
}

func colorText(text string, bashColorCode string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", bashColorCode, text)
}

func getLogTag(logLevel *types.LogLevel) string {
	if logLevel == nil {
		return "         "
	}
	return colorText(fmt.Sprintf("[%s]", colorText(cJust(logLevel.Name, 7, " "), logLevel.BashColorCode)), "1")
}

func InternalLog(message string, logLevel *types.LogLevel) {
	fmt.Printf("%s %s\033[K\n", getLogTag(logLevel), message)
}

func TempInternalLog(message string, logLevel *types.LogLevel) {
	fmt.Printf("%s %s\033[K\r", getLogTag(logLevel), message)
}
