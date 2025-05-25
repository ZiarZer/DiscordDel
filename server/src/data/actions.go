package data

import (
	"errors"

	"github.com/ZiarZer/DiscordDel/types"
	"github.com/ZiarZer/DiscordDel/utils"
)

func (repo *Repository) createActionsTable() error {
	_, err := repo.db.Exec(
		"CREATE TABLE IF NOT EXISTS `actions` (\n" +
			"`id` integer NOT NULL,\n" +
			"`type` varchar(10) NOT NULL,\n" +
			"`scope` varchar(10) NOT NULL,\n" +
			"`target_id` varchar(20),\n" +
			"`description` varchar(40) NOT NULL,\n" +
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

func (repo *Repository) InsertAction(action *types.Action) (int64, error) {
	if action == nil || action.Type == nil || action.Scope == nil {
		utils.InternalLog("Cannot insert action because of invalid information", utils.ERROR)
		return 0, errors.New("inserting an action requires at least a type and a scope")
	}

	stmt, err := repo.db.Prepare("INSERT INTO `actions` (`type`, `scope`, `target_id`, `description`) VALUES (?, ?, ?, ?)")
	if err != nil {
		utils.InternalLog("Failed to prepare action insertion", utils.ERROR)
		return 0, err
	}
	result, err := stmt.Exec(*action.Type, *action.Scope, action.TargetId, action.Description)
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
