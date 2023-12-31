package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wildberries-tech/wblogger"

	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
)

type DbMetrics interface {
	ReadStatsFromDB(s *sql.DB)
}

type SqlMetrics interface {
	Inc(query, success string)
	WriteTiming(start time.Time, query, success string)
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type Txs map[byte]Tx

var _ Tx = (pgx.Tx)(nil)

type Repository struct {
	db      *pgxpool.Pool
	sqlDB   *sql.DB
	metrics SqlMetrics
	timeout time.Duration
}

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SchemaName      string
	MaxOpenConns    int
	MaxIdleConns    int
	MinConns        int
	MaxConnLifetime time.Duration
	Timeout         time.Duration
}

const (
	defaultMapLen      = 1000
	defaultSelectLimit = 10

	metricsSuccess = "true"
	metricsFail    = "false"

	sql_Where = " WHERE "
	sql_And   = " AND "
	sql_Or    = " OR "
	sql_Asc   = " ASC"
	sql_Desc  = " DESC"
)

func NewRepository(cfg Config, dbMetrics DbMetrics, sqlMetrics SqlMetrics) (*Repository, error) {
	ctx := context.Background()
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		wblogger.Error(ctx, "NewPgRepo-ParseConfig", err)
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}

	sqlDB, err := sql.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}

	//config.ConnConfig.PreferSimpleProtocol = true
	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MinConns = int32(cfg.MinConns)
	if cfg.MaxConnLifetime > 0 {
		config.MaxConnLifetime = cfg.MaxConnLifetime
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		wblogger.Error(ctx, "NewPgRepo-ConnectConfig", err)
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}

	if err = db.Ping(ctx); err != nil {
		wblogger.Error(ctx, "NewPgRepo-Ping", err)
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	go func(m DbMetrics, updatePeriod time.Duration, ctx context.Context) {
		ticker := time.NewTicker(updatePeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
			// Безопасно для закрытой БД
			m.ReadStatsFromDB(sqlDB)
		}
	}(dbMetrics, 5*time.Second, ctx)

	return &Repository{
		db:      db,
		sqlDB:   sqlDB,
		metrics: sqlMetrics,
		timeout: timeout,
	}, nil
}

func (r *Repository) Close() {
	r.db.Close()
	r.sqlDB.Close()
}

func (r *Repository) SqlDB() *sql.DB {
	return r.sqlDB
}

// Begin используется для создания транзакции и её дальнейшей передачи в методы стора
func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	const metricName = "Begin"

	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	start := time.Now().UTC()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		wblogger.Error(ctx, "Begin-BeginTx-err", err)
		return nil, errors.Wrap(apperror.ErrInternal, "Begin transaction error: "+err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return tx, nil
}

func (r *Repository) BeginWithOptions(ctx context.Context, opts *pgx.TxOptions) (Tx, error) {
	const metricName = "BeginWithOptions"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	start := time.Now().UTC()

	tx, err := r.db.BeginTx(ctx, *opts)
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		wblogger.Error(ctx, "Begin-BeginTx-err", err)
		return nil, errors.Wrap(apperror.ErrInternal, "Begin transaction error: "+err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return tx, nil
}

func (r *Repository) Exec(ctx context.Context, sql string, arguments ...interface{}) error {
	const metricName = "Exec"
	_, err := r.db.Exec(ctx, sql, arguments...)
	if err != nil {
		wblogger.Error(ctx, metricName+" error", err)
		return errors.Wrap(apperror.ErrInternal, metricName+" error: "+err.Error())
	}
	return nil
}

func (r *Repository) ExecTx(ctx context.Context, tx Tx, sql string, arguments ...interface{}) error {
	const metricName = "ExecTx"
	_, err := tx.Exec(ctx, sql, arguments...)
	if err != nil {
		wblogger.Error(ctx, metricName+" error", err)
		return errors.Wrap(apperror.ErrInternal, metricName+" error: "+err.Error())
	}
	return nil
}
