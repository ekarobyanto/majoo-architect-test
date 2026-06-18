package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestHealthRepository_Ping(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewHealthRepository(sqlxDB)

	t.Run("success", func(t *testing.T) {
		mock.ExpectPing()
		err := repo.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(errors.New("db error"))
		err := repo.Ping(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
