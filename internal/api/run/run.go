// Package run provides application initialization and startup logic.
package run

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"

	"open-replays/internal/api/config"
	sqlite "open-replays/internal/api/db"
	"open-replays/internal/api/http/router"
	v1 "open-replays/internal/api/http/v1"
	"open-replays/internal/api/logger"
	repodb "open-replays/internal/api/repository/sqlite"
	"open-replays/internal/api/repository/storage"
	"open-replays/internal/api/usecase"
	"open-replays/internal/api/validation"
)

// Run initializes and starts the API server.
func Run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	}, os.Stdout)
	if err != nil {
		return err
	}
	slog.SetDefault(log)

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

	videoValidator := validation.NewVideoValidator()

	localStorage := storage.NewLocalStorage(cfg.Storage.BaseDir, cfg.Storage.PublicURL)
	videosRepo := repodb.NewVideosRepo(db)
	thumbnailsService := usecase.NewMetadataService(localStorage)
	videoProcessor := usecase.NewVideoProcessor(
		thumbnailsService,
		videosRepo,
		localStorage,
		cfg.Video.WorkerCount,
	)
	videosService := usecase.NewVideosService(
		videosRepo,
		localStorage,
		videoProcessor,
	)
	videosHandler := v1.NewVideosHandler(videosService, videoValidator)

	router.Register(r, videosHandler)

	slog.Info("starting server", "address", cfg.Server.Address())

	return r.Run(cfg.Server.Address())
}
