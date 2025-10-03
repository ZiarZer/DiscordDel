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
			"`object_id` varchar(20) NOT NULL,\n" +
			"`object_type` varchar(7) NOT NULL,\n" +
			"`oldest_read_id` varchar(20) NOT NULL,\n" +
			"`newest_read_id` varchar(20) NOT NULL,\n" +
			"`reached_top` boolean NOT NULL DEFAULT FALSE,\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`object_id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `crawling` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) GetCrawlingInfo(objectId types.Snowflake) (*types.CrawlingInfo, error) {
	stmt, err := repo.db.Prepare("SELECT `object_id`, `oldest_read_id`, `newest_read_id`, `reached_top` FROM `crawling` WHERE `object_id` = ?")
	if err != nil {
		println(err.Error())
		utils.InternalLog("Failed to prepare getting channel crawling info", utils.ERROR)
		return nil, err
	}
	var crawlingInfo types.CrawlingInfo
	err = stmt.QueryRow(objectId).Scan(&crawlingInfo.ObjectId, &crawlingInfo.OldestReadId, &crawlingInfo.NewestReadId, &crawlingInfo.ReachedTop)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		utils.InternalLog("Failed to get channel crawling info", utils.ERROR)
	}
	return &crawlingInfo, nil
}

func (repo *Repository) InsertCrawlingInfo(objectId types.Snowflake, objectType string, oldestReadId types.Snowflake, newestReadId types.Snowflake, reachedTop bool) error {
	stmt, err := repo.db.Prepare("INSERT INTO `crawling` (`object_id`, `object_type`, `oldest_read_id`, `newest_read_id`, `reached_top`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		utils.InternalLog("Failed to prepare new crawling info insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(objectId, objectType, oldestReadId, newestReadId, reachedTop)

	if err != nil {
		utils.InternalLog("Failed to insert new crawling info", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) UpdateCrawlingInfo(objectId types.Snowflake, oldestReadId *types.Snowflake, newestReadId *types.Snowflake, reachedTop *bool) error {
	queryParams := []interface{}{}
	preparedUpdates := []string{}
	if oldestReadId != nil {
		queryParams = append(queryParams, *oldestReadId)
		preparedUpdates = append(preparedUpdates, "`oldest_read_id` = ?")
	}
	if newestReadId != nil {
		queryParams = append(queryParams, *newestReadId)
		preparedUpdates = append(preparedUpdates, "`newest_read_id` = ?")
	}
	if reachedTop != nil {
		queryParams = append(queryParams, *reachedTop)
		preparedUpdates = append(preparedUpdates, "`reached_top` = ?")
	}
	if len(queryParams) == 0 {
		return nil
	}
	queryParams = append(queryParams, objectId)

	query := fmt.Sprintf("UPDATE `crawling` SET %s WHERE `object_id` = ?", strings.Join(preparedUpdates, ", "))
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		utils.InternalLog("Failed to prepare new crawling info insertion", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(queryParams...)

	if err != nil {
		utils.InternalLog("Failed to update crawling info", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) UpdateCrawlingOldestReadId(objectId types.Snowflake, updatedOldestReadId types.Snowflake) error {
	stmt, err := repo.db.Prepare("UPDATE `crawling` SET `oldest_read_id` = ? WHERE `object_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare crawling info update", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(updatedOldestReadId, objectId)

	if err != nil {
		utils.InternalLog("Failed to update crawling info", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) UpdateCrawlingNewestReadId(objectId types.Snowflake, updatedNewestReadId types.Snowflake) error {
	stmt, err := repo.db.Prepare("UPDATE `crawling` SET `newest_read_id` = ? WHERE `object_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare crawling info update", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(updatedNewestReadId, objectId)

	if err != nil {
		utils.InternalLog("Failed to update crawling info", utils.ERROR)
		return err
	}
	return nil
}

func (repo *Repository) UpdateCrawlingReachedTop(objectId types.Snowflake, updatedReachedTop bool) error {
	stmt, err := repo.db.Prepare("UPDATE `crawling` SET `reached_top` = ? WHERE `object_id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare crawling info update", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(updatedReachedTop, objectId)

	if err != nil {
		utils.InternalLog("Failed to update crawling info", utils.ERROR)
		return err
	}
	return nil
}
