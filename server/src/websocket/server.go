package websocket

import (
	"fmt"
	"net/http"

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
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			utils.Log("Failed to read WebSocket message from client", utils.ERROR)
			break
		}
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			utils.Log("Failed to send WebSocket message to client", utils.ERROR)
		}
	}
}

func RunWebSocketServer(pattern string, port int) {
	http.HandleFunc(pattern, handleConnection)
	utils.Log(fmt.Sprintf("Websocket server started: ws://localhost:%d", port), utils.INFO)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
