package data

import (
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createGuildsTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `guilds` (\n" +
			"`id` varchar(20) NOT NULL,\n" +
			"`name` varchar(50) NOT NULL,\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `guilds` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) InsertGuild(guild types.Guild) error {
	stmt, err := repo.db.Prepare("INSERT INTO `guilds` (`id`, `name`) VALUES(?, ?) ON CONFLICT DO NOTHING")
	if err != nil {
		utils.InternalLog("Failed to prepare guild insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(guild.Id, guild.Name)
	if err != nil {
		utils.InternalLog("Failed to insert guild", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) InsertMultipleGuilds(guilds []types.Guild) error {
	if len(guilds) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `guilds` (`id`, `name`) VALUES %s ON CONFLICT DO NOTHING",
		strings.TrimSuffix(strings.Repeat("(?, ?), ", len(guilds)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare guilds insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(guilds)*2)
	for i := range guilds {
		params[2*i] = guilds[i].Id
		params[2*i+1] = guilds[i].Name
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog("Failed to insert guilds", utils.ERROR)
		return err
	}
	return nil
}
