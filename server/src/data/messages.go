package data

import (
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createMessagesTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `messages` (\n" +
			"`id` varchar(20) NOT NULL,\n" +
			"`content` varchar(5000) NOT NULL,\n" +
			"`type` int NOT NULL,\n" +
			"`channel_id` varchar(20) NOT NULL,\n" +
			"`author_id` varchar(20) NOT NULL,\n" +
			"`pinned` boolean NOT NULL DEFAULT FALSE,\n" +
			"`status` varchar(10) NOT NULL DEFAULT 'PENDING',\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `messages` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) InsertMultipleMessages(messages []types.Message) error {
	if len(messages) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `messages` (`id`, `content`, `type`, `channel_id`, `author_id`, `pinned`) VALUES %s\n"+
			"ON CONFLICT DO UPDATE SET `content` = EXCLUDED.`content`, `type` = EXCLUDED.`type`, `pinned` = EXCLUDED.`pinned`",
		strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?, ?, ?), ", len(messages)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare messages insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(messages)*6)
	for i := range messages {
		params[6*i] = messages[i].Id
		params[6*i+1] = messages[i].Content
		params[6*i+2] = messages[i].Type
		params[6*i+3] = messages[i].ChannelId
		params[6*i+4] = messages[i].Author.Id
		if messages[i].InteractionMetadata != nil {
			params[6*i+4] = messages[i].InteractionMetadata.Triggerer.Id
		}
		params[6*i+5] = messages[i].Pinned
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog(fmt.Sprintf("Failed to insert messages: %s", err.Error()), utils.ERROR)
		return err
	}
	return nil
}
