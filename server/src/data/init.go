package data

import (
	"database/sql"

	"github.com/ZiarZer/DiscordDel/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func (repo *Repository) initTables() {
	err := repo.createGuildsTable()
	if err != nil {
		panic(err)
	}
	err = repo.createChannelsTable()
	if err != nil {
		panic(err)
	}
	err = repo.createMessagesTable()
	if err != nil {
		panic(err)
	}
	err = repo.createCrawlingTable()
	if err != nil {
		panic(err)
	}
}

func NewRepository() *Repository {
	db, err := sql.Open("sqlite3", ".discorddel.db")
	if err != nil {
		utils.InternalLog("Failed to open database file", utils.FATAL)
		panic(err)
	}
	repo := Repository{db: db}
	repo.initTables()
	return &repo
}
