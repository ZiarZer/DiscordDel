package wsserver

import (
	"fmt"
	"net/http"

	"github.com/ZiarZer/DiscordDel/crawl"
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

	for {
		err := handleMessage(conn)
		if err != nil {
			utils.InternalLog("Closing server WebSocket", utils.INFO)
			return
		}
	}
}

var sdk discord.DiscordSdk
var crawler crawl.Crawler

func RunWebSocketServer(pattern string, port int) {
	http.HandleFunc(pattern, handleConnection)
	sdk = discord.DiscordSdk{}
	crawler = crawl.Crawler{Sdk: &sdk}
	utils.InternalLog(fmt.Sprintf("Websocket server started: ws://localhost:%d", port), utils.INFO)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
