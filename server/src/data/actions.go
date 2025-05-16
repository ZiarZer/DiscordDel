package data

import (
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createActionsTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `actions` (\n" +
			"`id` integer NOT NULL,\n" +
			"`title` varchar(40) NOT NULL,\n" +
			"`start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"`end_time` timestamp,\n" +
			"`end_reason` varchar(40),\n" +
			"`status` varchar(10) NOT NULL DEFAULT 'RUNNING',\n" +
			"`last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"PRIMARY KEY(`id`))",
	)
	if err != nil {
		utils.InternalLog("Failed to create `actions` table", utils.FATAL)
		return err
	}
	return nil
}

func (repo *Repository) InsertAction(title string) (int64, error) {
	stmt, err := repo.db.Prepare("INSERT INTO `actions` (`title`) VALUES (?)")
	if err != nil {
		utils.InternalLog("Failed to prepare action insertion", utils.ERROR)
		return 0, err
	}
	result, err := stmt.Exec(title)
	if err != nil {
		utils.InternalLog("Failed to insert action", utils.ERROR)
		return 0, err
	}
	return result.LastInsertId()
}

func (repo *Repository) EndAction(actionId int64, endReason string) error {
	stmt, err := repo.db.Prepare("UPDATE `actions` SET `end_time` = CURRENT_TIMESTAMP, `status` = 'ENDED', `end_reason` = ? WHERE `id` = ?")
	if err != nil {
		utils.InternalLog("Failed to prepare action update", utils.ERROR)
		return err
	}
	_, err = stmt.Exec(endReason, actionId)
	if err != nil {
		utils.InternalLog("Failed to update action", utils.ERROR)
		return err
	}
	return nil
}
