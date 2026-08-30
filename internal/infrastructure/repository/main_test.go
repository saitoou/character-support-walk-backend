package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/chocoko/character-support-walk-backend/internal/test/testutils"
	"github.com/stretchr/testify/require"
)

var repositoryTestDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	testContainer, err := testutils.StartPostgres(ctx)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to setup repository test database: %v\n",
			err,
		)
		os.Exit(1)
	}

	repositoryTestDB = testContainer.DB

	exitCode := m.Run()

	if err := testContainer.Close(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to clean up repository test database: %v\n",
			err,
		)

		if exitCode == 0 {
			exitCode = 1
		}
	}

	os.Exit(exitCode)

}

func beginTestTx(t *testing.T) *sql.Tx {
	t.Helper()

	tx, err := repositoryTestDB.Begin()
	require.NoError(t, err)

	t.Cleanup(func() {
		err := tx.Rollback()

		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Errorf(
				"failed to rollback test transaction: %v\n",
				err,
			)
		}
	})

	return tx
}
