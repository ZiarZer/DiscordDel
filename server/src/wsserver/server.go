package wsserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ZiarZer/DiscordDel/actions"
	"github.com/ZiarZer/DiscordDel/crawl"
	"github.com/ZiarZer/DiscordDel/data"
	"github.com/ZiarZer/DiscordDel/delete"
	"github.com/ZiarZer/DiscordDel/discord"
	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.InternalLog("Failed to upgrade connection to WebSocket", utils.ERROR)
	}
	sdk.SetWsConn(conn)
	sdk.Log("Connected to client via WebSocket", utils.INFO)
	defer conn.Close()

	ctx := context.Background()
	for {
		err := handleMessage(ctx, conn)
		if err != nil {
			currentActionMutex.Lock()
			defer currentActionMutex.Unlock()
			if currentAction {
				cancelCurrentAction()
			}
			return
		}
	}
}

var sdk discord.DiscordSdk
var crawler crawl.Crawler
var deleter delete.Deleter
var actionLogger actions.ActionLogger

func RunWebSocketServer(pattern string, port int) {
	http.HandleFunc(pattern, handleConnection)
	sdk = discord.DiscordSdk{Repo: data.NewRepository(), ApiClient: &discord.ApiClient{Delay: 700}}
	actionLogger = actions.ActionLogger{Repo: sdk.Repo}
	crawler = crawl.Crawler{Sdk: &sdk, ActionLogger: &actionLogger}
	deleter = delete.Deleter{Sdk: &sdk, ActionLogger: &actionLogger}
	utils.InternalLog(fmt.Sprintf("Websocket server started: ws://localhost:%d", port), utils.INFO)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
