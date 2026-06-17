package database

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/gomega"
)

func TestTransactor_WithinTransaction(t *testing.T) {
	g := NewWithT(t)

	dbRaw, mock, err := sqlmock.New()
	g.Expect(err).NotTo(HaveOccurred())
	defer dbRaw.Close()

	db := sqlx.NewDb(dbRaw, "postgres")
	transactor := NewTransactor(db)
	ctx := context.Background()

	t.Run("should commit transaction on success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			// Verify tx is in context
			tx := txCtx.Value(txKey)
			g.Expect(tx).To(BeAssignableToTypeOf(&sqlx.Tx{}))
			return nil
		})

		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(mock.ExpectationsWereMet()).To(Succeed())
	})

	t.Run("should rollback transaction on error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		appErr := errors.New("app error")
		err := transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			return appErr
		})

		g.Expect(err).To(Equal(appErr))
		g.Expect(mock.ExpectationsWereMet()).To(Succeed())
	})

	t.Run("should rollback transaction on panic", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		defer func() {
			p := recover()
			g.Expect(p).To(Equal("panic error"))
			g.Expect(mock.ExpectationsWereMet()).To(Succeed())
		}()

		_ = transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			panic("panic error")
		})
	})
}

func TestGetQueryer(t *testing.T) {
	g := NewWithT(t)

	dbRaw, _, _ := sqlmock.New()
	db := sqlx.NewDb(dbRaw, "postgres")
	ctx := context.Background()

	t.Run("should return DB if no transaction in context", func(t *testing.T) {
		queryer := GetQueryer(ctx, db)
		g.Expect(queryer).To(Equal(db))
	})

	t.Run("should return Tx if transaction in context", func(t *testing.T) {
		dbRaw2, mock2, _ := sqlmock.New()
		db2 := sqlx.NewDb(dbRaw2, "postgres")
		mock2.ExpectBegin()
		tx, _ := db2.BeginTxx(ctx, nil)

		txCtx := context.WithValue(ctx, txKey, tx)
		queryer := GetQueryer(txCtx, db)
		g.Expect(queryer).To(Equal(tx))
	})
}
