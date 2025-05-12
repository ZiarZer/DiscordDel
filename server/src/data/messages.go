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

func (repo *Repository) InsertMultipleMessages(messages []types.Message, status string) error {
	if len(messages) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO `messages` (`id`, `content`, `type`, `channel_id`, `author_id`, `pinned`, `status`) VALUES %s\n"+
			"ON CONFLICT DO UPDATE SET `content` = EXCLUDED.`content`, `type` = EXCLUDED.`type`, `pinned` = EXCLUDED.`pinned`",
		strings.TrimSuffix(strings.Repeat("(?, ?, ?, ?, ?, ?, ?), ", len(messages)), ", "),
	)
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare messages insertion", utils.ERROR)
		return err
	}

	params := make([]interface{}, len(messages)*7)
	for i := range messages {
		params[7*i] = messages[i].Id
		params[7*i+1] = messages[i].Content
		params[7*i+2] = messages[i].Type
		params[7*i+3] = messages[i].ChannelId
		params[7*i+4] = messages[i].Author.Id
		if messages[i].InteractionMetadata != nil {
			params[7*i+4] = messages[i].InteractionMetadata.Triggerer.Id
		}
		params[7*i+5] = messages[i].Pinned
		params[7*i+6] = status
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		utils.InternalLog(fmt.Sprintf("Failed to insert messages: %s", err.Error()), utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) GetMessagesByChannelId(channelId types.Snowflake, authorIds []types.Snowflake) ([]types.Message, error) {
	if len(authorIds) == 0 {
		return []types.Message{}, nil
	}
	authorIdsParams := strings.TrimSuffix(strings.Repeat("?, ", len(authorIds)), ", ")
	stmt, err := repo.db.Prepare(
		fmt.Sprintf(
			"SELECT `id`, `content`, `type`, `channel_id`, `author_id`, `pinned`, `status` FROM `messages` WHERE `channel_id` = ? AND `author_id` IN (%s) ORDER BY `id`",
			authorIdsParams,
		),
	)
	if err != nil {
		utils.InternalLog("Failed to prepare getting messages by channel ID", utils.ERROR)
		return nil, err
	}
	params := []any{types.Snowflake(channelId)}
	for i := range authorIds {
		params = append(params, string(authorIds[i]))
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		utils.InternalLog("Failed to get messages by channel ID", utils.ERROR)
		return nil, err
	}
	defer rows.Close()
	var messages []types.Message
	for rows.Next() {
		var message types.Message
		err = rows.Scan(&message.Id, &message.Content, &message.Type, &message.ChannelId, &message.Author.Id, &message.Pinned, &message.Status)
		messages = append(messages, message)
		if err != nil {
			return messages, err
		}
	}
	return messages, nil
}

func (repo *Repository) UpdateMessagesStatus(messageIds []types.Snowflake, updatedStatus string) error {
	if len(messageIds) == 0 {
		return nil
	}
	messageIdsParams := strings.TrimSuffix(strings.Repeat("?, ", len(messageIds)), ", ")
	stmt, err := repo.db.Prepare(fmt.Sprintf("UPDATE `messages` SET `status` = ? WHERE `id` IN (%s)", messageIdsParams))
	if err != nil {
		utils.InternalLog("Failed to prepare message status update", utils.ERROR)
		return err
	}
	params := []any{updatedStatus}
	for i := range messageIds {
		params = append(params, messageIds[i])
	}
	_, err = stmt.Exec(params...)

	if err != nil {
		utils.InternalLog("Failed to update message status", utils.ERROR)
		return err
	}
	return nil
}
