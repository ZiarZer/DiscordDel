package main

import (
	"fmt"

	"github.com/ZiarZer/DiscordDel/utils"
)

const Version = "0.0.1"

func main() {
	utils.Log(fmt.Sprintf("DiscordDel - v%s", Version), nil)
}
