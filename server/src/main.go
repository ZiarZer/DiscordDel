package main

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/utils"
	"github.com/ZiarZer/DiscordDel/wsserver"
)

const Version = "0.9.1"

func main() {
	utils.InternalLog(fmt.Sprintf("DiscordDel - v%s", Version), nil)
	wsserver.RunWebSocketServer("/", 8765)
}
