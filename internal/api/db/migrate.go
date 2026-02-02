package db

import (
	"context"
	"database/sql"

	_ "embed"
)

//go:embed schema/videos.sql
var _videosDDL string

// Migrate runs database migrations.
func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, _videosDDL)
	return err
}
