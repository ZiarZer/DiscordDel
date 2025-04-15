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

type RequestBody interface {
	handle(conn *websocket.Conn) error
}

type LoginRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

type GetGuildRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
	GuildId            string `json:"guildId"`
}

var bodyConstructors = map[string]func() RequestBody{
	"LOGIN":     func() RequestBody { return &LoginRequestBody{} },
	"GET_GUILD": func() RequestBody { return &GetGuildRequestBody{} },
}

func handleMessage(conn *websocket.Conn) error {
	_, stringMessage, err := conn.ReadMessage()
	if err != nil {
		utils.Log("Client disconnected", utils.INFO)
		return err
	}

	var message Message
	err = json.Unmarshal(stringMessage, &message)
	if err != nil {
		utils.Log("Failed to parse WebSocket message from client", utils.FATAL)
		return err
	}

	body := bodyConstructors[message.Type]()
	err = json.Unmarshal(message.Body, body)
	if err != nil {
		utils.Log("Failed to read Websocket message's body", utils.FATAL)
		return err
	}
	body.handle(conn)
	return nil
}

func (body *LoginRequestBody) handle(conn *websocket.Conn) error {
	user := discord.Login(body.AuthorizationToken)
	jsonUser, err := json.Marshal(user)
	if err != nil {
		utils.Log("Failed to serialize user info", utils.ERROR)
		return err
	}
	response := Message{
		Type: "LOGIN",
		Body: json.RawMessage(jsonUser),
	}
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize LOGIN response", utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}

func (body *GetGuildRequestBody) handle(conn *websocket.Conn) error {
	guild := discord.GetGuild(body.GuildId, body.AuthorizationToken)
	jsonGuild, err := json.Marshal(guild)
	if err != nil {
		utils.Log("Failed to serialize guild info", utils.ERROR)
		return err
	}
	response := Message{
		Type: "GET_GUILD",
		Body: json.RawMessage(jsonGuild),
	}
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize GET_GUILD response", utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}
