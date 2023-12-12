package v1

import (
	t "github.com/golanguzb70/simple-post-app/api/tokens"
	"github.com/golanguzb70/simple-post-app/config"
	"github.com/golanguzb70/simple-post-app/pkg/logger"
	"github.com/golanguzb70/simple-post-app/storage"
	"github.com/golanguzb70/simple-post-app/storage/redisrepo"
)

type handlerV1 struct {
	log        *logger.Logger
	cfg        config.Config
	storage    storage.StorageI
	jwthandler t.JWTHandler
	redis      redisrepo.InMemoryStorageI
}

type HandlerV1Config struct {
	Logger     *logger.Logger
	Cfg        config.Config
	Postgres   storage.StorageI
	JWTHandler t.JWTHandler
	Redis      redisrepo.InMemoryStorageI
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:        c.Logger,
		cfg:        c.Cfg,
		storage:    c.Postgres,
		jwthandler: c.JWTHandler,
		redis:      c.Redis,
	}
}
