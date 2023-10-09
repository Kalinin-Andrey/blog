package redis

import (
	"context"
	"time"

	"github.com/minipkg/db/redis"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
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

func (r *Repository) Lock(ctx context.Context, key string) (bool, error) {
	const metricName = "Lock"
	err := r.db.Client().Watch(ctx, func(tx *goredis.Tx) error {
		// Get current value or zero.
		t, err := tx.Get(ctx, key).Time()
		if err != nil && err != goredis.Nil {
			return err
		}

		t = time.Now().UTC()
		// Operation is committed only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe goredis.Pipeliner) error {
			pipe.Set(ctx, key, t.Format(time.RFC3339Nano), 0)
			return nil
		})

		if err != nil {
			return err
		}
		return err
	}, key)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, goredis.TxFailedErr) {
		return false, nil
	}
	return false, errors.Wrap(apperror.ErrInternal, metricName+" error:"+err.Error())
}
