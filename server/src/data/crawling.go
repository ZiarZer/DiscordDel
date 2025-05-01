package data

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createCrawlingTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `crawling` (\n" +
			"`channel_id` varchar(20) NOT NULL,\n" +
			"`oldest_read_message_id` varchar(20) NOT NULL,\n" +
			"`newest_read_message_id` varchar(20) NOT NULL,\n" +
			"`reached_top` boolean NOT NULL DEFAULT FALSE,\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`channel_id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `crawling` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) GetChannelCrawlingInfo(channelId string) (*types.CrawlingInfo, error) {
	stmt, err := repo.db.Prepare("SELECT `channel_id`, `oldest_read_message_id`, `newest_read_message_id`, `reached_top` FROM `crawling` WHERE `channel_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to get channel crawling info", utils.ERROR)
		return nil, err
	}
	var crawlingInfo types.CrawlingInfo
	err = stmt.QueryRow(channelId).Scan(&crawlingInfo.ChannelId, &crawlingInfo.OldestReadMessageId, &crawlingInfo.NewestReadMessageId, &crawlingInfo.ReachedTop)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &crawlingInfo, nil
}

func (repo *Repository) InsertChannelCrawlingInfo(channelId string, oldestReadMessageId string, newestReadMessageId string, reachedTop bool) error {
	stmt, err := repo.db.Prepare("INSERT INTO `crawling` (`channel_id`, `oldest_read_message_id`, `newest_read_message_id`, `reached_top`) VALUES (?, ?, ?, ?)")
	if err != nil {
		utils.InternalLog("Failed to prepare new channel crawling info insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(channelId, oldestReadMessageId, newestReadMessageId, reachedTop)

	if err != nil {
		utils.InternalLog("Failed to insert new channel crawling info", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) UpdateChannelCrawlingInfo(channelId string, oldestReadMessageId *string, newestReadMessageId *string, reachedTop *bool) error {
	queryParams := []interface{}{}
	preparedUpdates := []string{}
	if oldestReadMessageId != nil {
		queryParams = append(queryParams, *oldestReadMessageId)
		preparedUpdates = append(preparedUpdates, "`oldest_read_message_id` = ?")
	}
	if newestReadMessageId != nil {
		queryParams = append(queryParams, *newestReadMessageId)
		preparedUpdates = append(preparedUpdates, "`newest_read_message_id` = ?")
	}
	if reachedTop != nil {
		queryParams = append(queryParams, *reachedTop)
		preparedUpdates = append(preparedUpdates, "`reached_top` = ?")
	}
	if len(queryParams) == 0 {
		return nil
	}
	queryParams = append(queryParams, channelId)

	query := fmt.Sprintf("UPDATE `crawling` SET %s WHERE `channel_id` = ?", strings.Join(preparedUpdates, ", "))
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare new channel crawling info insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(queryParams...)

	if err != nil {
		utils.InternalLog("Failed to update channel crawling info", utils.ERROR)
		return err
	}
	return nil
}
