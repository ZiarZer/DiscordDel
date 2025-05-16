package actions

import (
	"fmt"
	"time"

	"github.com/ZiarZer/DiscordDel/data"
	"github.com/ZiarZer/DiscordDel/utils"
)

type ActionLogger struct {
	Repo *data.Repository
}

type Action struct {
	Id         int64
	Title      string
	StartTime  time.Time
	LogFunc    func(message string, logLevel *utils.LogLevel)
	LogEndTime bool
	Major      bool
}

func (actionLogger *ActionLogger) StartAction(title string, logFunc func(message string, logLevel *utils.LogLevel), logEndTime bool, major bool) *Action {
	logFunc(title, utils.INFO)
	action := &Action{Title: title, StartTime: time.Now(), LogFunc: logFunc, LogEndTime: logEndTime, Major: major}
	if major {
		action.Id, _ = actionLogger.Repo.InsertAction(title)
	}
	return action
}

func (actionLogger *ActionLogger) EndAction(action *Action) {
	durationInSeconds := time.Now().Unix() - action.StartTime.Unix()
	if action.LogEndTime {
		action.LogFunc(fmt.Sprintf("[%s] finished in %s", action.Title, utils.FormatDuration(durationInSeconds)), utils.INFO)
	}
	if action.Major {
		actionLogger.Repo.EndAction(action.Id, "FINISHED")
	}
}
