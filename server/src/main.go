package main

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/ZiarZer/DiscordDel/websocket"
)

const Version = "0.1.0"

func main() {
	utils.Log(fmt.Sprintf("DiscordDel - v%s", Version), nil)
	websocket.RunWebSocketServer("/", 8765)
}
