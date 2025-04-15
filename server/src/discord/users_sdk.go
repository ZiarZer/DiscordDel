package discord

import (
	"encoding/json"
	"io"

	"github.com/ZiarZer/DiscordDel/utils"
)

func Login(authorizationToken string) *User {
	resp, err := login(authorizationToken)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log(err.Error(), utils.ERROR)
		return nil
	}
	if resp.StatusCode != 200 {
		utils.Log(string(body), utils.ERROR)
		return nil
	}
	var loggedUser User
	json.Unmarshal(body, &loggedUser)
	return &loggedUser
}
