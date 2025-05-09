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
	LogLevel string `json:"logLevel"`
	Message  string `json:"message"`
}

func (sdk *DiscordSdk) Log(message string, logLevel *utils.LogLevel) {
	utils.InternalLog(message, logLevel)

	if sdk.wsConn != nil {
		jsonLogEntry, err := json.Marshal(LogEntry{LogLevel: logLevel.Name, Message: message})
		if err != nil {
			utils.InternalLog("Failed to serialize log data", utils.ERROR)
			return
		}
		wsbase.SendMessageToClient(sdk.wsConn, "LOG", jsonLogEntry)
	}
}
