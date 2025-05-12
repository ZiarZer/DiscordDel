package actions

import (
	"fmt"
	"time"

	"github.com/ZiarZer/DiscordDel/utils"
)

type Action struct {
	Title      string
	StartTime  time.Time
	LogFunc    func(message string, logLevel *utils.LogLevel)
	LogEndTime bool
}

func StartAction(title string, logFunc func(message string, logLevel *utils.LogLevel), logEndTime bool) *Action {
	logFunc(title, utils.INFO)
	return &Action{Title: title, StartTime: time.Now(), LogFunc: logFunc, LogEndTime: logEndTime}
}

func (action *Action) EndAction() {
	durationInSeconds := time.Now().Unix() - action.StartTime.Unix()
	if action.LogEndTime {
		action.LogFunc(fmt.Sprintf("[%s] finished in %s", action.Title, utils.FormatDuration(durationInSeconds)), utils.INFO)
	}
}
