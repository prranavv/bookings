package dbrepo

import (
	"database/sql"

	"github.com/prranavv/bookings/internal/config"
	"github.com/prranavv/bookings/internal/models"
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

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}
