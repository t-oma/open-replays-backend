// Package run provides application initialization and startup logic.
package run

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"

	sqlite "open-replays/internal/api/db"
	"open-replays/internal/api/http/router"
	repodb "open-replays/internal/api/repository/sqlite"
	"open-replays/internal/api/repository/storage"
	"open-replays/internal/api/usecase"
)

// Run initializes and starts the API server.
func Run() error {
	ctx := context.Background()

	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "db.sqlite3"
	}

	// busy_timeout
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL", dbPath)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	if err = sqlite.Migrate(ctx, db); err != nil {
		return err
	}

	r := gin.Default()
	r.Use(cors.Default())

	localStorage := storage.NewLocalStorage("uploads", "http://localhost:8080/media")
	videosRepo := repodb.NewVideosRepo(db)
	thumbnailsService := usecase.NewMetadataService(localStorage)
	videoProcessor := usecase.NewVideoProcessor(thumbnailsService, videosRepo, localStorage, 2)
	videosUC := usecase.NewVideosService(videosRepo, localStorage, videoProcessor)

	router.Register(r, videosUC)

	return r.Run(":8080")
}
