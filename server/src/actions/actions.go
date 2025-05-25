package actions

import (
	"fmt"
	"time"

	"github.com/ZiarZer/DiscordDel/data"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

type ActionLogger struct {
	Repo *data.Repository
}

func NewMajorAction(actionType types.ActionType, scope types.Scope, targetId *types.Snowflake, description string) *types.Action {
	return &types.Action{
		Type:        &actionType,
		Scope:       &scope,
		TargetId:    targetId,
		Description: description,
		StartTime:   time.Now(),
	}
}

func NewAction(description string) *types.Action {
	return &types.Action{
		Description: description,
		StartTime:   time.Now(),
	}
}

func (actionLogger *ActionLogger) StartAction(action *types.Action, logFunc func(message string, logLevel *types.LogLevel), logEndTime bool) *types.Action {
	action.LogFunc = logFunc
	action.LogEndTime = logEndTime
	logFunc(action.Description, utils.INFO)
	if action.Type != nil && action.Scope != nil {
		insertId, _ := actionLogger.Repo.InsertAction(action)
		action.Id = &insertId
	}
	return action
}

func (actionLogger *ActionLogger) EndAction(action *types.Action) {
	durationInSeconds := time.Now().Unix() - action.StartTime.Unix()
	if action.LogEndTime {
		action.LogFunc(fmt.Sprintf("[%s] finished in %s", action.Description, utils.FormatDuration(durationInSeconds)), utils.INFO)
	}
	if action.Id != nil {
		actionLogger.Repo.EndAction(*action.Id, "FINISHED")
	}
}
