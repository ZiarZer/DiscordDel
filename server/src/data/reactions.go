package data

import (
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createReactionsTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `reactions` (\n" +
			"`message_id` varchar(20) NOT NULL,\n" +
			"`user_id` varchar(20) NOT NULL,\n" +
			"`emoji` varchar(20) NOT NULL,\n" +
			"`is_burst` boolean NOT NULL,\n" +
			"`status` varchar(10) NOT NULL DEFAULT 'PENDING',\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`message_id`, `user_id`, `emoji`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `reactions` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) InsertMultipleReactions(messageId types.Snowflake, userIds []types.Snowflake, emoji string, isBurst bool) error {
	if len(userIds) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `reactions` (`message_id`, `user_id`, `emoji`, `is_burst`) VALUES %s\n"+
			"ON CONFLICT DO UPDATE SET `is_burst` = EXCLUDED.`is_burst`",
		strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?), ", len(userIds)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare reactions insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(userIds)*4)
	for i := range userIds {
		params[4*i] = messageId
		params[4*i+1] = userIds[i]
		params[4*i+2] = emoji
		params[4*i+3] = isBurst
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog(fmt.Sprintf("Failed to insert reactions: %s", err.Error()), utils.ERROR)
		return err
	}
	return nil
}
