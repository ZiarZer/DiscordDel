package websocket

import (
	"encoding/json"

	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type string          `json:"type"`
	Body json.RawMessage `json:"body"`
}

type LoginBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

func handleMessage(conn *websocket.Conn) {
	messageType, stringMessage, err := conn.ReadMessage()
	if err != nil {
		utils.Log("Failed to read WebSocket message from client", utils.FATAL)
		return
	}

	var message Message
	err = json.Unmarshal(stringMessage, &message)
	if err != nil {
		utils.Log("Failed to parse WebSocket message from client", utils.FATAL)
		return
	}

	body := &LoginBody{}
	err = json.Unmarshal(message.Body, body)
	if err != nil {
		utils.Log("Failed to read Websocket message's body", utils.FATAL)
	}

	user := discord.Login(body.AuthorizationToken)
	jsonUser, err := json.Marshal(user)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
	}
	response := Message{
		Type: "LOGIN",
		Body: json.RawMessage(jsonUser),
	}
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize LOGIN response", utils.FATAL)
	}
	err = conn.WriteMessage(messageType, stringResponse)

	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
	}
}
