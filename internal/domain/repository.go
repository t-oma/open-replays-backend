package domain

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	mutex sync.Mutex
	db    *pgxpool.Pool
}
