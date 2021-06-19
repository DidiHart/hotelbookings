package dbrepo

import (
	"database/sql"

	"github.com/DidiHart/hotelbookings/internal/config"
	"github.com/DidiHart/hotelbookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
