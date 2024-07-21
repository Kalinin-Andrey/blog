package tsdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"blog/internal/pkg/apperror"
)

type DbMetrics interface {
	ReadStatsFromDB(s *sql.DB)
}

type SqlMetrics interface {
	Inc(query, success string)
	WriteTiming(start time.Time, query, success string)
}

type GaugeMetrics interface {
	Add(valueName string, value float64)
	Set(valueName string, value float64)
}

type CounterMetrics interface {
	Inc(labelValues ...string)
	Add(val int64, labelValues ...string)
}

type RepositoryMetrics struct {
	SqlMetrics     SqlMetrics
	GaugeMetrics   GaugeMetrics
	CounterMetrics CounterMetrics
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
	metrics *RepositoryMetrics
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
	metricsSuccess = "true"
	metricsFail    = "false"

	defaultCapacityForResult = 100
	errMsg_duplicateKey      = "duplicate key"

	sql_Where = " WHERE "
	sql_And   = " AND "
	sql_Or    = " OR "
	sql_Asc   = " ASC"
	sql_Desc  = " DESC"

	sql_OnConflictDoNothing = " ON CONFLICT DO NOTHING;"
)

func NewRepository(cfg Config, dbMetrics DbMetrics, metrics *RepositoryMetrics) (*Repository, error) {
	ctx := context.Background()
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("[%w] pgxpool.ParseConfig(url) error: %w", apperror.ErrInternal, err)
	}

	sqlDB, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("[%w] sql.Open() error: %w", apperror.ErrInternal, err)
	}

	//config.ConnConfig.PreferSimpleProtocol = true
	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MinConns = int32(cfg.MinConns)
	if cfg.MaxConnLifetime > 0 {
		config.MaxConnLifetime = cfg.MaxConnLifetime
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("[%w] pgxpool.NewWithConfig() error: %w", apperror.ErrInternal, err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("[%w] db.Ping() error: %w", apperror.ErrInternal, err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	go func(ctx context.Context, m DbMetrics, updatePeriod time.Duration) {
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
	}(ctx, dbMetrics, 5*time.Second)

	return &Repository{
		db:      db,
		sqlDB:   sqlDB,
		metrics: metrics,
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
	const metricName = "Repository.Begin"

	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	start := time.Now().UTC()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s begin transaction error: %w", apperror.ErrInternal, metricName, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return tx, nil
}
func (r *Repository) BeginWithOptions(ctx context.Context, opts *pgx.TxOptions) (Tx, error) {
	const metricName = "Repository.BeginWithOptions"

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	start := time.Now().UTC()

	tx, err := r.db.BeginTx(ctx, *opts)
	if err != nil {
		r.metrics.SqlMetrics.Inc(metricName, metricsFail)
		r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] %s begin transaction error: %w", apperror.ErrInternal, metricName, err)
	}
	r.metrics.SqlMetrics.Inc(metricName, metricsSuccess)
	r.metrics.SqlMetrics.WriteTiming(start, metricName, metricsSuccess)
	return tx, nil
}

func (r *Repository) Exec(ctx context.Context, sql string, arguments ...interface{}) error {
	const metricName = "Repository.Exec"
	_, err := r.db.Exec(ctx, sql, arguments...)
	if err != nil {
		return fmt.Errorf("[%w] %s SQL request exec error: %w", apperror.ErrInternal, metricName, err)
	}
	return nil
}

func (r *Repository) ExecTx(ctx context.Context, tx Tx, sql string, arguments ...interface{}) error {
	const metricName = "Repository.ExecTx"
	_, err := tx.Exec(ctx, sql, arguments...)
	if err != nil {
		return fmt.Errorf("[%w] %s SQL request exec tx error: %w", apperror.ErrInternal, metricName, err)
	}
	return nil
}
