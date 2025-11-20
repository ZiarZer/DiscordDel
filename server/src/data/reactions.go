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
			"`channel_id` varchar(20) NOT NULL,\n" +
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

func (repo *Repository) InsertMultipleReactions(channelId types.Snowflake, messageId types.Snowflake, userIds []types.Snowflake, emoji string, isBurst bool) error {
	if len(userIds) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `reactions` (`channel_id`, `message_id`, `user_id`, `emoji`, `is_burst`) VALUES %s\n"+
			"ON CONFLICT DO UPDATE SET `is_burst` = EXCLUDED.`is_burst`",
		strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?, ?), ", len(userIds)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare reactions insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(userIds)*5)
	for i := range userIds {
		params[5*i] = channelId
		params[5*i+1] = messageId
		params[5*i+2] = userIds[i]
		params[5*i+3] = emoji
		params[5*i+4] = isBurst
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog(fmt.Sprintf("Failed to insert reactions: %s", err.Error()), utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) GetPendingReactionsByChannelId(channelId types.Snowflake, authorIds []types.Snowflake) ([]types.Reaction, error) {
	if len(authorIds) == 0 {
		return []types.Reaction{}, nil
	}
	authorIdsParams := strings.TrimSuffix(strings.Repeat("?, ", len(authorIds)), ", ")
	stmt, err := repo.db.Prepare(
		fmt.Sprintf(
			"SELECT `channel_id`, `message_id`, `user_id`, `emoji`, `is_burst`, `status` FROM `reactions`\n"+
				"WHERE `channel_id` = ? AND `user_id` IN (%s) AND `status` = 'PENDING'\n"+
				"ORDER BY `message_id`, `emoji`, `user_id`",
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
		err = rows.Scan(&reaction.ChannelId, &reaction.MessageId, &reaction.UserId, &reaction.Emoji, &reaction.IsBurst, &reaction.Status)
		reactions = append(reactions, reaction)
		if err != nil {
			return reactions, err
		}
	}
	return reactions, nil
}

func (repo *Repository) UpdateReactionsStatus(reactions []types.Reaction, updatedStatus types.CrawlingStatus) error {
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

func (repo *Repository) UpdateReactionsStatusByMessageId(messageId types.Snowflake, updatedStatus types.CrawlingStatus) error {
	stmt, err := repo.db.Prepare("UPDATE `reactions` SET `status` = ? WHERE `message_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare reactions status update", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(updatedStatus, messageId)
	if err != nil {
		utils.InternalLog("Failed to update reactions status", utils.ERROR)
		return err
	}
	return nil
}
