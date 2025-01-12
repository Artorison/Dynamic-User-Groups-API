package database

import (
	"API/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDBConnection(cfg config.AppConfig) (*Database, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

func (d *Database) Close() error {
	return d.DB.Close()
}
