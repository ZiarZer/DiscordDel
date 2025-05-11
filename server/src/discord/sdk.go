package discord

import (
	"encoding/json"

	"github.com/ZiarZer/DiscordDel/data"

	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/ZiarZer/DiscordDel/wsbase"
	"github.com/gorilla/websocket"
)

type DiscordSdk struct {
	wsConn    *websocket.Conn
	Repo      *data.Repository
	ApiClient *ApiClient
}

func (sdk *DiscordSdk) SetWsConn(wsConn *websocket.Conn) {
	sdk.wsConn = wsConn
}

type LogEntry struct {
	LogLevel *string `json:"logLevel"`
	Message  string  `json:"message"`
}

func (sdk *DiscordSdk) sendLogEntryToClient(message string, logLevel *utils.LogLevel, websocketMessageType string) {
	if sdk.wsConn != nil {
		var logLevelName *string
		if logLevel != nil {
			logLevelName = &logLevel.Name
		}
		jsonLogEntry, err := json.Marshal(LogEntry{LogLevel: logLevelName, Message: message})
		if err != nil {
			utils.InternalLog("Failed to serialize log data", utils.ERROR)
			return
		}
		wsbase.SendMessageToClient(sdk.wsConn, websocketMessageType, jsonLogEntry)
	}
}

func (sdk *DiscordSdk) TempLog(message string, logLevel *utils.LogLevel) {
	utils.TempInternalLog(message, logLevel)

	sdk.sendLogEntryToClient(message, logLevel, "TEMP_LOG")
}

func (sdk *DiscordSdk) Log(message string, logLevel *utils.LogLevel) {
	utils.InternalLog(message, logLevel)

	sdk.sendLogEntryToClient(message, logLevel, "LOG")
}
