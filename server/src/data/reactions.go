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

func (repo *Repository) GetReactionsByChannelId(channelId types.Snowflake, authorIds []types.Snowflake) ([]types.Reaction, error) {
	if len(authorIds) == 0 {
		return []types.Reaction{}, nil
	}
	authorIdsParams := strings.TrimSuffix(strings.Repeat("?, ", len(authorIds)), ", ")
	stmt, err := repo.db.Prepare(
		fmt.Sprintf(
			"SELECT `r`.`message_id`, `r`.`user_id`, `r`.`emoji`, `r`.`is_burst`, `r`.`status`\n"+
				"FROM `reactions` AS `r` JOIN `messages` AS `m` ON `m`.`id` = `r`.`message_id`\n"+
				"WHERE `m`.`channel_id` = ? AND `r`.`user_id` IN (%s) ORDER BY `r`.`message_id`, `r`.`emoji`, `r`.`user_id`",
			authorIdsParams,
		),
	)
	if err != nil {
		utils.InternalLog("Failed to prepare getting reactions by channel ID", utils.ERROR)
		return nil, err
	}
	params := []any{types.Snowflake(channelId)}
	for i := range authorIds {
		params = append(params, string(authorIds[i]))
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		utils.InternalLog("Failed to get reactions by channel ID", utils.ERROR)
		return nil, err
	}
	defer rows.Close()
	var reactions []types.Reaction
	for rows.Next() {
		var reaction types.Reaction
		err = rows.Scan(&reaction.MessageId, &reaction.UserId, &reaction.Emoji, &reaction.IsBurst, &reaction.Status)
		reactions = append(reactions, reaction)
		if err != nil {
			return reactions, err
		}
	}
	return reactions, nil
}

func (repo *Repository) UpdateReactionsStatus(reactions []types.Reaction, updatedStatus string) error {
	if len(reactions) == 0 {
		return nil
	}
	messageIdsParams := strings.TrimSuffix(strings.Repeat("(?, ?, ?), ", len(reactions)), ", ")
	stmt, err := repo.db.Prepare(fmt.Sprintf("UPDATE `reactions` SET `status` = ? WHERE (`message_id`, `user_id`, `emoji`) IN (%s)", messageIdsParams))
	if err != nil {
		utils.InternalLog("Failed to prepare reactions status update", utils.ERROR)
		return err
	}
	params := []any{updatedStatus}
	for i := range reactions {
		params = append(params, reactions[i].MessageId, reactions[i].UserId, reactions[i].Emoji)
	}
	_, err = stmt.Exec(params...)

	if err != nil {
		utils.InternalLog("Failed to update reactions status", utils.ERROR)
		return err
	}
	return nil
}
