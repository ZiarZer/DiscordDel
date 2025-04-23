package wsserver

import (
	"encoding/json"

	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/ZiarZer/DiscordDel/wsbase"
	"github.com/gorilla/websocket"
)

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
		utils.InternalLog("Client disconnected", utils.INFO)
		return err
	}

	var message wsbase.Message
	err = json.Unmarshal(stringMessage, &message)
	if err != nil {
		utils.InternalLog("Failed to parse WebSocket message from client", utils.FATAL)
		return err
	}

	body := bodyConstructors[message.Type]()
	err = json.Unmarshal(message.Body, body)
	if err != nil {
		utils.InternalLog("Failed to read Websocket message's body", utils.FATAL)
		return err
	}
	body.handle(conn)
	return nil
}

func (body *LoginRequestBody) handle(conn *websocket.Conn) error {
	user := sdk.Login(body.AuthorizationToken)
	jsonUser, err := json.Marshal(user)
	if err != nil {
		utils.InternalLog("Failed to serialize user info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "LOGIN", jsonUser)
}

func (body *GetGuildRequestBody) handle(conn *websocket.Conn) error {
	guild := sdk.GetGuild(body.GuildId, body.AuthorizationToken)
	jsonGuild, err := json.Marshal(guild)
	if err != nil {
		utils.InternalLog("Failed to serialize guild info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD", jsonGuild)
}

func (body *GetChannelRequestBody) handle(conn *websocket.Conn) error {
	channel := sdk.GetChannel(body.ChannelId, body.AuthorizationToken)
	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		utils.InternalLog("Failed to serialize channel info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_CHANNEL", jsonChannel)
}

func (body *GetUserGuildsRequestBody) handle(conn *websocket.Conn) error {
	guilds := sdk.GetUserGuilds(body.AuthorizationToken)
	jsonGuildList, err := json.Marshal(guilds)
	if err != nil {
		utils.InternalLog("Failed to serialize guilds list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_USER_GUILDS", jsonGuildList)
}

func (body *GetGuildChannelsRequestBody) handle(conn *websocket.Conn) error {
	channels := sdk.GetGuildChannels(body.GuildId, body.AuthorizationToken)
	jsonChannelList, err := json.Marshal(channels)
	if err != nil {
		utils.InternalLog("Failed to serialize channels list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD_CHANNELS", jsonChannelList)
}
