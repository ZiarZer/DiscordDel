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

type ActionType string

const (
	CRAWL  ActionType = "CRAWL"
	DELETE ActionType = "DELETE"
)

type Scope string

const (
	CHANNEL Scope = "CHANNEL"
	GUILD   Scope = "GUILD"
	ALL     Scope = "ALL"
)

type StartActionRequestBody struct {
	AuthorizationToken string                `json:"authorizationToken"`
	AuthorIds          []types.Snowflake     `json:"authorIds"`
	Type               ActionType            `json:"type"`
	Scope              Scope                 `json:"scope"`
	TargetId           *types.Snowflake      `json:"targetId"`
	Options            *delete.DeleteOptions `json:"options"`
}

var bodyConstructors = map[string]func() RequestBody{
	"LOGIN":              func() RequestBody { return &LoginRequestBody{} },
	"GET_USER_GUILDS":    func() RequestBody { return &GetUserGuildsRequestBody{} },
	"GET_GUILD":          func() RequestBody { return &GetGuildRequestBody{} },
	"GET_GUILD_CHANNELS": func() RequestBody { return &GetGuildChannelsRequestBody{} },
	"GET_CHANNEL":        func() RequestBody { return &GetChannelRequestBody{} },
	"START_ACTION":       func() RequestBody { return &StartActionRequestBody{} },
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
func endAction() {
	currentActionMutex.Lock()
	defer currentActionMutex.Unlock()
	currentAction = false
}

func (body *LoginRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	user, err := sdk.Login(authorizedContext)
	if err != nil {
		return err
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		utils.InternalLog("Failed to serialize user info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "LOGIN", jsonUser)
}

func (body *GetGuildRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	guild, err := sdk.GetGuild(authorizedContext, body.GuildId)
	if err != nil {
		return err
	}
	jsonGuild, err := json.Marshal(guild)
	if err != nil {
		utils.InternalLog("Failed to serialize guild info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD", jsonGuild)
}

func (body *GetChannelRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	channel, err := sdk.GetChannel(authorizedContext, body.ChannelId)
	if err != nil {
		return err
	}
	jsonChannel, err := json.Marshal(channel)
	if err != nil {
		utils.InternalLog("Failed to serialize channel info", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_CHANNEL", jsonChannel)
}

func (body *GetUserGuildsRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	guilds, err := sdk.GetUserGuilds(authorizedContext)
	if err != nil {
		return err
	}
	jsonGuildList, err := json.Marshal(guilds)
	if err != nil {
		utils.InternalLog("Failed to serialize guilds list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_USER_GUILDS", jsonGuildList)
}

func (body *GetGuildChannelsRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	channels, err := sdk.GetGuildChannels(authorizedContext, body.GuildId)
	if err != nil {
		return err
	}
	jsonChannelList, err := json.Marshal(channels)
	if err != nil {
		utils.InternalLog("Failed to serialize channels list", utils.ERROR)
		return err
	}
	return wsbase.SendMessageToClient(conn, "GET_GUILD_CHANNELS", jsonChannelList)
}

func (body *StartActionRequestBody) handle(ctx context.Context, conn *websocket.Conn) error {
	if body.Type != CRAWL && body.Type != DELETE {
		sdk.Log("Unknown action type", utils.ERROR)
		return nil
	} else if body.Scope != CHANNEL || body.Scope != GUILD && body.Scope == ALL {
		sdk.Log("Unknown action scope", utils.ERROR)
		return nil
	} else if (body.Scope == CHANNEL || body.Scope == GUILD) && body.TargetId == nil {
		sdk.Log("No target ID specified for action", utils.ERROR)
		return nil
	} else if len(body.AuthorIds) == 0 {
		sdk.Log("No author IDs specified for action", utils.ERROR)
		return nil
	}

	if !startAction() {
		return nil
	}
	defer endAction()

	authorizedContext := context.WithValue(ctx, types.CtxKey{Key: "authorizationToken"}, body.AuthorizationToken)
	if body.Type == CRAWL {
		if body.Scope == CHANNEL {
			crawler.CrawlChannel(authorizedContext, body.AuthorIds, *body.TargetId)
		} else if body.Scope == GUILD {
			crawler.CrawlGuild(authorizedContext, body.AuthorIds, *body.TargetId)
		} else if body.Scope == ALL {
			crawler.CrawlAllGuilds(authorizedContext, body.AuthorIds)
		}
	} else if body.Type == DELETE {
		var options delete.DeleteOptions
		if body.Options != nil {
			options = *body.Options
		}
		if body.Scope == CHANNEL {
			deleter.DeleteChannelCrawledData(authorizedContext, body.AuthorIds, *body.TargetId, options)
		} else if body.Scope == GUILD {
			deleter.BulkDeleteCrawledData(authorizedContext, body.AuthorIds, body.TargetId, options)
		} else if body.Scope == ALL {
			deleter.BulkDeleteCrawledData(authorizedContext, body.AuthorIds, nil, options)
		}
	}
	return nil
}
