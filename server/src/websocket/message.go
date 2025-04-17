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

func newMessage(Type string, jsonBytes []byte) Message {
	return Message{
		Type: Type,
		Body: json.RawMessage(jsonBytes),
	}
}

type RequestBody interface {
	handle(conn *websocket.Conn) error
}

type LoginRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

type GetUserGuildsRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

type GetGuildRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
	GuildId            string `json:"guildId"`
}

type GetGuildChannelsRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
	GuildId            string `json:"guildId"`
}

type GetChannelRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
	ChannelId          string `json:"channelId"`
}

var bodyConstructors = map[string]func() RequestBody{
	"LOGIN":              func() RequestBody { return &LoginRequestBody{} },
	"GET_USER_GUILDS":    func() RequestBody { return &GetUserGuildsRequestBody{} },
	"GET_GUILD":          func() RequestBody { return &GetGuildRequestBody{} },
	"GET_GUILD_CHANNELS": func() RequestBody { return &GetGuildChannelsRequestBody{} },
	"GET_CHANNEL":        func() RequestBody { return &GetChannelRequestBody{} },
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
	response := newMessage("LOGIN", jsonUser)
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
	response := newMessage("GET_GUILD", jsonGuild)
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

func (body *GetChannelRequestBody) handle(conn *websocket.Conn) error {
	channel := discord.GetChannel(body.ChannelId, body.AuthorizationToken)
	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		utils.Log("Failed to serialize channel info", utils.ERROR)
		return err
	}
	response := newMessage("GET_CHANNEL", jsonChannel)
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize GET_CHANNEL response", utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}

func (body *GetUserGuildsRequestBody) handle(conn *websocket.Conn) error {
	guilds := discord.GetUserGuilds(body.AuthorizationToken)
	jsonGuildList, err := json.Marshal(guilds)
	if err != nil {
		utils.Log("Failed to serialize guilds list", utils.ERROR)
		return err
	}
	response := newMessage("GET_USER_GUILDS", jsonGuildList)
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize GET_USER_GUILDS response", utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}

func (body *GetGuildChannelsRequestBody) handle(conn *websocket.Conn) error {
	channels := discord.GetGuildChannels(body.GuildId, body.AuthorizationToken)
	jsonChannelList, err := json.Marshal(channels)
	if err != nil {
		utils.Log("Failed to serialize channels list", utils.ERROR)
		return err
	}
	response := newMessage("GET_GUILD_CHANNELS", jsonChannelList)
	stringResponse, err := json.Marshal(response)
	if err != nil {
		utils.Log("Failed to serialize GET_GUILD_CHANNELS response", utils.FATAL)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, stringResponse)
	if err != nil {
		utils.Log("Failed to send WebSocket message to client", utils.FATAL)
		return err
	}
	return nil
}
