package run

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	sqlite "open-replays/api/internal/api/db"
	"open-replays/api/internal/api/http/router"
	repodb "open-replays/api/internal/api/repository/sqlite"
	"open-replays/api/internal/api/repository/storage"
	"open-replays/api/internal/api/usecase"
)

func Run() error {
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
	defer db.Close()

	if err := sqlite.Migrate(db); err != nil {
		return err
	}

	r := gin.Default()
	r.Use(cors.Default())

	storage := storage.NewLocalStorage("uploads/videos")
	videosRepo := repodb.NewVideosRepo(db)
	videosUC := usecase.NewVideosService(videosRepo, storage)

	router.Register(r, videosUC)

	return r.Run(":8080")
}
