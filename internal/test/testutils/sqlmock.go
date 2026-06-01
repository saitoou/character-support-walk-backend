package testutils

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db, mock
}
