// Package run provides application initialization and startup logic.
package run

import (
	"context"
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"

	"open-replays/internal/api/config"
	sqlite "open-replays/internal/api/db"
	"open-replays/internal/api/http/router"
	repodb "open-replays/internal/api/repository/sqlite"
	"open-replays/internal/api/repository/storage"
	"open-replays/internal/api/usecase"
)

// Run initializes and starts the API server.
func Run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	maxFileSize, err := cfg.Video.MaxFileSizeBytes()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", cfg.Database.DSN())
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

	localStorage := storage.NewLocalStorage(cfg.Storage.BaseDir, cfg.Storage.PublicURL)
	videosRepo := repodb.NewVideosRepo(db)
	thumbnailsService := usecase.NewMetadataService(localStorage)
	videoProcessor := usecase.NewVideoProcessor(
		thumbnailsService,
		videosRepo,
		localStorage,
		cfg.Video.WorkerCount,
	)
	videosUC := usecase.NewVideosService(
		videosRepo,
		localStorage,
		videoProcessor,
		maxFileSize,
		cfg.Video.AllowedExtensions,
	)

	router.Register(r, videosUC)

	return r.Run(":8080")
}
