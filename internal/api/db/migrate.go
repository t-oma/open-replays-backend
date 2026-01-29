package db

import (
	"database/sql"
	_ "embed"
)

//go:embed schema/videos.sql
var videosDDL string

// Migrate runs database migrations.
func Migrate(db *sql.DB) error {
	_, err := db.Exec(videosDDL)
	return err
}
