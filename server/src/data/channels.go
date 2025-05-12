package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createChannelsTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `channels` (\n" +
			"`id` varchar(20) NOT NULL,\n" +
			"`name` varchar(50),\n" +
			"`type` int NOT NULL,\n" +
			"`guild_id` varchar(20),\n" +
			"`parent_id` varchar(20),\n" +
			"`last_message_id` varchar(20),\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `channels` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) InsertChannel(channel types.Channel) error {
	stmt, err := repo.db.Prepare(
		"INSERT INTO `channels` (`id`, `name`, `type`, `guild_id`, `parent_id`, `last_message_id`) VALUES (?, ?, ?, ?, ?, ?)\n" +
			"ON CONFLICT DO UPDATE SET `name` = EXCLUDED.`name`, `type` = EXCLUDED.`type`, `parent_id` = EXCLUDED.`parent_id`, `last_message_id` = EXCLUDED.`last_message_id`",
	)
	if err != nil {
		utils.InternalLog("Failed to prepare channel insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(channel.Id, channel.Name, channel.Type, channel.GuildId, channel.ParentId, channel.LastMessageId)
	if err != nil {
		utils.InternalLog("Failed to insert channel", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) InsertMultipleChannels(channels []types.Channel) error {
	if len(channels) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `channels` (`id`, `name`, `type`, `guild_id`, `parent_id`, `last_message_id`) VALUES %s\n"+
			"ON CONFLICT DO UPDATE SET `name` = EXCLUDED.`name`, `type` = EXCLUDED.`type`, `parent_id` = EXCLUDED.`parent_id`, `last_message_id` = EXCLUDED.`last_message_id`",
		strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?, ?, ?), ", len(channels)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare channels insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(channels)*6)
	for i := range channels {
		params[6*i] = channels[i].Id
		params[6*i+1] = channels[i].Name
		params[6*i+2] = channels[i].Type
		params[6*i+3] = channels[i].GuildId
		params[6*i+4] = channels[i].ParentId
		params[6*i+5] = channels[i].LastMessageId
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog(fmt.Sprintf("Failed to insert channels: %s", err.Error()), utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) GetChannelChildrenCount(parentChannelId types.Snowflake) (int, error) {
	stmt, err := repo.db.Prepare("SELECT count(*) FROM `channels` WHERE `parent_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare getting channel children count", utils.ERROR)
		return 0, err
	}
	row := stmt.QueryRow(parentChannelId)

	var count int
	err = row.Scan(&count)
	if err != nil {
		utils.InternalLog("Failed to scan children count", utils.ERROR)
		return 0, err
	}
	return count, nil
}

func (repo *Repository) GetChannelsWithPendingMessages(authorIds []types.Snowflake, guildId *types.Snowflake) ([]types.Snowflake, error) {
	query := "SELECT DISTINCT `channel_id` FROM `messages` WHERE `status` = 'PENDING'"
	var rows *sql.Rows
	var err error
	if guildId != nil {
		query += " AND `guild_id` = ?"
		stmt, err := repo.db.Prepare(query)
		if err != nil {
			utils.InternalLog("Failed to prepare getting guild channels with PENDING messages", utils.ERROR)
			return nil, err
		}
		rows, err = stmt.Query(guildId)
	} else {
		rows, err = repo.db.Query(query)
	}

	if err != nil {
		utils.InternalLog("Failed to get channels with PENDING messages", utils.ERROR)
		return nil, err
	}
	defer rows.Close()
	var channelIds []types.Snowflake
	for rows.Next() {
		var channelId types.Snowflake
		err := rows.Scan(&channelId)
		channelIds = append(channelIds, channelId)
		if err != nil {
			return channelIds, err
		}
	}
	return channelIds, nil
}
