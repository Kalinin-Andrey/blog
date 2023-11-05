package redis

import (
	"time"

	"github.com/minipkg/db/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	RedisNil = "redis: nil"

	metricsSuccess = "true"
	metricsFail    = "false"
)

type RedisMetrics interface {
	Inc(query, success string)
	WriteTiming(startTime time.Time, query, success string)
}

type Repository struct {
	db      redis.IDB
	metrics RedisMetrics
}

func NewRepository(cfg redis.Config, metrics RedisMetrics) (*Repository, error) {
	db, err := redis.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:      db,
		metrics: metrics,
	}, nil
}

func (r *Repository) DB() goredis.Cmdable {
	return r.db.DB()
}

func (r *Repository) Close() error {
	return r.db.Close()
}
