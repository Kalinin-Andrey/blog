package infrastructure

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minipkg/prometheus-utils"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/wildberries-tech/wblogger"

	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/redis"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb_cluster"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
	"github.com/Kalinin-Andrey/blog/migrations"
	_ "github.com/Kalinin-Andrey/blog/migrations"
)

const (
	PgMigrationLockKey = "pgMigrate"
)

type AppConfig struct {
	NameSpace   string
	Name        string
	Service     string
	Environment string
}

type Infrastructure struct {
	appConfig *AppConfig
	config    *Config
	TsDB      *tsdb_cluster.ReplicaSet
	Redis     *redis.ReplicaSet
}

func New(ctx context.Context, appConfig *AppConfig, config *Config) (*Infrastructure, error) {
	infra := &Infrastructure{
		appConfig: appConfig,
		config:    config,
	}

	wblogger.Info(ctx, "Redis: try to connect...")
	if err := infra.redisInit(ctx); err != nil {
		return nil, err
	}
	wblogger.Info(ctx, "Redis: connected.")

	wblogger.Info(ctx, "TsDB try to connect...")
	if err := infra.tsDBInit(ctx); err != nil {
		return nil, err
	}
	wblogger.Info(ctx, "TsDB: connected.")

	return infra, nil
}

func (infra *Infrastructure) tsDBInit(ctx context.Context) error {
	pgConf := infra.config.TsDB.getConfig()
	nodeLabel := "tsdb_master"
	sqlMetrics := prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[0].DbName)
	dbMetrics := prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[0].DbName)
	gaugeMetrics := prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[0].DbName)
	master, err := tsdb.NewRepository(pgConf[0], dbMetrics, sqlMetrics, gaugeMetrics)
	if err != nil {
		return err
	}

	nodeLabel = "tsdb_replicas"
	sqlMetrics = prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[1].DbName)
	dbMetrics = prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[1].DbName)
	gaugeMetrics = prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, pgConf[1].DbName)
	slave, err := tsdb.NewRepository(pgConf[1], dbMetrics, sqlMetrics, gaugeMetrics)
	if err != nil {
		return err
	}
	wblogger.Info(ctx, "TsDB connected.")
	infra.TsDB = tsdb_cluster.NewReplicaSet(master, slave)

	if err := infra.tsDBMigrate(ctx); err != nil {
		return errors.Wrap(apperror.ErrInternal, "app.Infra.tsDBMigrate() error: "+err.Error())
	}

	return nil
}

func (infra *Infrastructure) tsDBMigrate(ctx context.Context) (err error) {
	wblogger.Info(ctx, "TsDBMigration start")

	repo := infra.TsDB.Master()
	migrations.CurrentRepo = repo

	goose.SetBaseFS(migrations.EmbedMigrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	if err := goose.Up(repo.SqlDB(), "."); err != nil {
		return err
	}
	wblogger.Info(ctx, "TsDBMigration done!")
	return nil
}

func (infra *Infrastructure) redisInit(ctx context.Context) (err error) {
	conf := infra.config.Redis.getConfig()
	nodeLabel := "redis_master"
	metrics := prometheus_utils.NewRedisMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel)
	master, err := redis.NewRepository(conf[0], metrics)
	if err != nil {
		return err
	}

	nodeLabel = "redis_replicas"
	metrics = prometheus_utils.NewRedisMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel)
	slave, err := redis.NewRepository(conf[1], metrics)
	if err != nil {
		return err
	}
	wblogger.Info(ctx, "redis connected.")
	infra.Redis = redis.NewReplicaSet(master, slave)

	return err
}

func (infra *Infrastructure) Close() error {
	infra.TsDB.Close()
	infra.Redis.Close()
	return nil
}
