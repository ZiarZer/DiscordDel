package wsserver

import (
	"fmt"
	"net/http"

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
		utils.Log("Failed to upgrade connection to WebSocket", utils.ERROR)
	}
	utils.Log("Connected to client via WebSocket", utils.INFO)
	defer conn.Close()

	for {
		err := handleMessage(conn)
		if err != nil {
			utils.Log("Closing server WebSocket", utils.INFO)
			return
		}
	}
}

var sdk discord.DiscordSdk

func RunWebSocketServer(pattern string, port int) {
	http.HandleFunc(pattern, handleConnection)
	sdk = discord.DiscordSdk{}
	utils.Log(fmt.Sprintf("Websocket server started: ws://localhost:%d", port), utils.INFO)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
