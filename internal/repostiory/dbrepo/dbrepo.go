package dbrepo

import (
	"database/sql"

	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/repostiory"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repostiory.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewtestingRepo(a *config.AppConfig) repostiory.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
