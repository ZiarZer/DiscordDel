package wsserver

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/ZiarZer/DiscordDel/delete"
	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/ZiarZer/DiscordDel/wsbase"
	"github.com/gorilla/websocket"
)

func makeAuthorizationTokenContext(parentCtx context.Context, authorizationToken string) context.Context {
	return context.WithValue(parentCtx, types.CtxKey{Key: "authorizationToken"}, authorizationToken)
}

type RequestBody interface {
	handle(ctx context.Context, conn *websocket.Conn) error
}

type LoginRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

type GetUserGuildsRequestBody struct {
	AuthorizationToken string `json:"authorizationToken"`
}

type GetGuildRequestBody struct {
	AuthorizationToken string          `json:"authorizationToken"`
	GuildId            types.Snowflake `json:"guildId"`
}

type GetGuildChannelsRequestBody struct {
	AuthorizationToken string          `json:"authorizationToken"`
	GuildId            types.Snowflake `json:"guildId"`
}

type GetChannelRequestBody struct {
	AuthorizationToken string          `json:"authorizationToken"`
	ChannelId          types.Snowflake `json:"channelId"`
}

type CrawlChannelRequestBody struct {
	AuthorizationToken string            `json:"authorizationToken"`
	ChannelId          types.Snowflake   `json:"channelId"`
	AuthorIds          []types.Snowflake `json:"authorIds"`
}

type CrawlGuildRequestBody struct {
	AuthorizationToken string            `json:"authorizationToken"`
	GuildId            types.Snowflake   `json:"guildId"`
	AuthorIds          []types.Snowflake `json:"authorIds"`
}

type CrawlAllGuildsRequestBody struct {
	AuthorizationToken string            `json:"authorizationToken"`
	AuthorIds          []types.Snowflake `json:"authorIds"`
}

type DeleteChannelDataRequestBody struct {
	AuthorizationToken string               `json:"authorizationToken"`
	AuthorIds          []types.Snowflake    `json:"authorIds"`
	ChannelId          types.Snowflake      `json:"channelId"`
	Options            delete.DeleteOptions `json:"options"`
}

type DeleteGuildDataRequestBody struct {
	AuthorizationToken string               `json:"authorizationToken"`
	AuthorIds          []types.Snowflake    `json:"authorIds"`
	GuildId            types.Snowflake      `json:"guildId"`
	Options            delete.DeleteOptions `json:"options"`
}

type DeleteAllDataRequestBody struct {
	AuthorizationToken string               `json:"authorizationToken"`
	AuthorIds          []types.Snowflake    `json:"authorIds"`
	Options            delete.DeleteOptions `json:"options"`
}

var bodyConstructors = map[string]func() RequestBody{
	"LOGIN":               func() RequestBody { return &LoginRequestBody{} },
	"GET_USER_GUILDS":     func() RequestBody { return &GetUserGuildsRequestBody{} },
	"GET_GUILD":           func() RequestBody { return &GetGuildRequestBody{} },
	"GET_GUILD_CHANNELS":  func() RequestBody { return &GetGuildChannelsRequestBody{} },
	"GET_CHANNEL":         func() RequestBody { return &GetChannelRequestBody{} },
	"CRAWL_CHANNEL":       func() RequestBody { return &CrawlChannelRequestBody{} },
	"CRAWL_GUILD":         func() RequestBody { return &CrawlGuildRequestBody{} },
	"CRAWL_ALL_GUILDS":    func() RequestBody { return &CrawlAllGuildsRequestBody{} },
	"DELETE_CHANNEL_DATA": func() RequestBody { return &DeleteChannelDataRequestBody{} },
	"DELETE_GUILD_DATA":   func() RequestBody { return &DeleteGuildDataRequestBody{} },
	"DELETE_ALL_DATA":     func() RequestBody { return &DeleteAllDataRequestBody{} },
}

func handleMessage(ctx context.Context, conn *websocket.Conn) error {
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
	go body.handle(context.Context(ctx), conn)
	return nil
}

var currentAction bool
var currentActionMutex sync.Mutex

func startAction() bool {
	currentActionMutex.Lock()
	defer currentActionMutex.Unlock()
	if currentAction {
		sdk.Log("An action is already running", utils.ERROR)
		return false
	}
	currentAction = true
	return true
}
func endAction() { currentAction = false }

func (body *LoginRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	user := sdk.Login(makeAuthorizationTokenContext(ctx, body.AuthorizationToken))
	jsonUser, err := json.Marshal(user)
	if err != nil {
		utils.InternalLog("Failed to serialize user info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "LOGIN", jsonUser)
}

func (body *GetGuildRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	guild := sdk.GetGuild(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.GuildId)
	jsonGuild, err := json.Marshal(guild)
	if err != nil {
		utils.InternalLog("Failed to serialize guild info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD", jsonGuild)
}

func (body *GetChannelRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	channel := sdk.GetChannel(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.ChannelId)
	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		utils.InternalLog("Failed to serialize channel info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_CHANNEL", jsonChannel)
}

func (body *GetUserGuildsRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	guilds := sdk.GetUserGuilds(makeAuthorizationTokenContext(ctx, body.AuthorizationToken))
	jsonGuildList, err := json.Marshal(guilds)
	if err != nil {
		utils.InternalLog("Failed to serialize guilds list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_USER_GUILDS", jsonGuildList)
}

func (body *GetGuildChannelsRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	channels := sdk.GetGuildChannels(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.GuildId)
	jsonChannelList, err := json.Marshal(channels)
	if err != nil {
		utils.InternalLog("Failed to serialize channels list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD_CHANNELS", jsonChannelList)
}

func (body *CrawlChannelRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()

	crawler.CrawlChannel(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds, body.ChannelId)
	return nil
}

func (body *CrawlGuildRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()
	crawler.CrawlGuild(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds, body.GuildId)
	return nil
}

func (body *CrawlAllGuildsRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()
	crawler.CrawlAllGuilds(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds)
	return nil
}

func (body *DeleteChannelDataRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()
	deleter.DeleteChannelCrawledData(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds, body.ChannelId, body.Options)
	return nil
}

func (body *DeleteGuildDataRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()
	deleter.BulkDeleteCrawledData(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds, &body.GuildId, body.Options)
	return nil
}

func (body *DeleteAllDataRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if !startAction() {
		return nil
	}
	defer endAction()
	deleter.BulkDeleteCrawledData(makeAuthorizationTokenContext(ctx, body.AuthorizationToken), body.AuthorIds, nil, body.Options)
	return nil
}
