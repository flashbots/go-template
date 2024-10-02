package database

import (
	"os"
	"testing"

	"github.com/flashbots/go-template/common"
	"github.com/flashbots/go-template/database/migrations"
	"github.com/flashbots/go-template/database/vars"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var (
	runDBTests = os.Getenv("RUN_DB_TESTS") == "1" //|| true
	testDBDSN  = common.GetEnv("TEST_DB_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
)

func resetDatabase(t *testing.T) *DatabaseService {
	t.Helper()
	if !runDBTests {
		t.Skip("Skipping database tests")
	}

	// Wipe test database
	_db, err := sqlx.Connect("postgres", testDBDSN)
	require.NoError(t, err)
	_, err = _db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`)
	require.NoError(t, err)

	db, err := NewDatabaseService(testDBDSN)
	require.NoError(t, err)
	return db
}

func TestMigrations(t *testing.T) {
	db := resetDatabase(t)
	query := `SELECT COUNT(*) FROM ` + vars.TableMigrations + `;`
	rowCount := 0
	err := db.DB.QueryRow(query).Scan(&rowCount)
	require.NoError(t, err)
	require.Len(t, migrations.Migrations.Migrations, rowCount)
}

func Test_DB1(t *testing.T) {
	db := resetDatabase(t)
	x, err := db.SomeQuery()
	require.NoError(t, err)
	require.Equal(t, uint64(0), x)
}
