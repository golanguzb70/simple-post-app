package postgres

import (
	"time"

	"github.com/golanguzb70/simple-post-app/config"
	"github.com/golanguzb70/simple-post-app/pkg/db"
	"github.com/golanguzb70/simple-post-app/pkg/logger"
)

var (
	CreatedAt time.Time
	UpdatedAt time.Time
)

type postgresRepo struct {
	Db  *db.Postgres
	Log *logger.Logger
	Cfg config.Config
}

func New(db *db.Postgres, log *logger.Logger, cfg config.Config) PostgresI {
	return &postgresRepo{
		Db:  db,
		Log: log,
		Cfg: cfg,
	}
}
