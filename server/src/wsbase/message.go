package wsbase

import (
	"encoding/json"
	"fmt"

	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

func newMessage(Type string, jsonBytes []byte) Message {
	return Message{
		Type: Type,
		Body: json.RawMessage(jsonBytes),
	}
}

func SendMessageToClient(conn *websocket.Conn, messageType string, messageBody []byte) error {
	response := newMessage(messageType, messageBody)
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log(fmt.Sprintf("Failed to serialize %s response", messageType), utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}
