// Package database exposes the postgres database
package database

import (
	"os"

	"github.com/flashbots/go-template/database/migrations"
	"github.com/flashbots/go-template/database/vars"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type DatabaseService struct {
	DB *sqlx.DB
}

func NewDatabaseService(dsn string) (*DatabaseService, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.DB.SetMaxOpenConns(50)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetConnMaxIdleTime(0)

	if os.Getenv("DB_DONT_APPLY_SCHEMA") == "" {
		migrate.SetTable(vars.TableMigrations)
		_, err := migrate.Exec(db.DB, "postgres", migrations.Migrations, migrate.Up)
		if err != nil {
			return nil, err
		}
	}

	dbService := &DatabaseService{DB: db} //nolint:exhaustruct
	err = dbService.prepareNamedQueries()
	return dbService, err
}

func (s *DatabaseService) prepareNamedQueries() (err error) {
	return nil
}

func (s *DatabaseService) Close() error {
	return s.DB.Close()
}

func (s *DatabaseService) SomeQuery() (count uint64, err error) {
	query := `SELECT COUNT(*) FROM ` + vars.TableTest + `;`
	row := s.DB.QueryRow(query)
	err = row.Scan(&count)
	return count, err
}
