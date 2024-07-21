package infrastructure

import (
	"blog/internal/pkg/appconst"
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minipkg/prometheus-utils"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"blog/internal/infrastructure/repository/redis"
	"blog/internal/infrastructure/repository/tsdb"
	"blog/internal/infrastructure/repository/tsdb_cluster"
	_ "blog/migrations/tsdb"
	tsdbmigration "blog/migrations/tsdb"
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
	Logger    *zap.Logger
	TsDB      *tsdb_cluster.ReplicaSet
	Redis     *redis.ReplicaSet
}

func New(ctx context.Context, appConfig *AppConfig, config *Config) (*Infrastructure, error) {
	infra := &Infrastructure{
		appConfig: appConfig,
		config:    config,
	}
	if err := infra.loggerInit(ctx, appConfig.Environment); err != nil {
		return nil, err
	}
	infra.Logger.Info("Logger initialised.")

	infra.Logger.Info("Redis: try to connect...")
	if err := infra.redisInit(ctx); err != nil {
		return nil, err
	}
	infra.Logger.Info("Redis: connected.")

	infra.Logger.Info("TsDB try to connect...")
	if err := infra.tsDBInit(ctx); err != nil {
		return nil, err
	}
	infra.Logger.Info("TsDB: connected.")

	return infra, nil
}

func (infra *Infrastructure) loggerInit(ctx context.Context, environment string) (err error) {
	var lvl zap.AtomicLevel

	switch environment {
	case appconst.Env_Local:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case appconst.Env_Dev:
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case appconst.Env_Stage:
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case appconst.Env_Prd:
		//lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	conf := zap.NewProductionConfig()
	conf.Level = lvl
	conf.OutputPaths = []string{"stdout"}
	conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	infra.Logger, err = conf.Build(zap.AddCallerSkip(1)) // записываем не logger в каждую запись, а реального вызывающего

	return err
}

func (infra *Infrastructure) tsDBInit(ctx context.Context) error {
	conf := infra.config.TsDB.getConfig()
	counterMetrics := prometheus_utils.NewCounter(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, "nb", "tsdb", "value_name", "success")

	nodeLabel := "tsdb_master"
	sqlMetrics := prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	dbMetrics := prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	gaugeMetrics := prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	master, err := tsdb.NewRepository(conf[0], dbMetrics, &tsdb.RepositoryMetrics{
		SqlMetrics:     sqlMetrics,
		GaugeMetrics:   gaugeMetrics,
		CounterMetrics: counterMetrics,
	})
	if err != nil {
		return err
	}

	nodeLabel = "tsdb_replicas"
	sqlMetrics = prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	dbMetrics = prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	gaugeMetrics = prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	slave, err := tsdb.NewRepository(conf[1], dbMetrics, &tsdb.RepositoryMetrics{
		SqlMetrics:     sqlMetrics,
		GaugeMetrics:   gaugeMetrics,
		CounterMetrics: counterMetrics,
	})
	if err != nil {
		return err
	}
	infra.Logger.Info("TsDB connected.")
	infra.TsDB = tsdb_cluster.NewReplicaSet(master, slave)

	if err := infra.tsDBMigrate(ctx); err != nil {
		return fmt.Errorf("app.Infra.tsDBMigrate() error: %w", err)
	}

	return nil
}

func (infra *Infrastructure) tsDBMigrate(ctx context.Context) (err error) {
	infra.Logger.Info("TsDBMigration start")

	repo := infra.TsDB.WriteRepo()
	tsdbmigration.CurrentRepo = repo

	goose.SetBaseFS(tsdbmigration.EmbedMigrations)
	if err := goose.SetDialect("pgx"); err != nil {
		return err
	}

	if err := goose.Up(repo.SqlDB(), "."); err != nil {
		return err
	}
	infra.Logger.Info("TsDBMigration done!")
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
	infra.Logger.Info("redis connected.")
	infra.Redis = redis.NewReplicaSet(master, slave)

	return err
}

func (infra *Infrastructure) Close() error {
	infra.TsDB.Close()
	infra.Redis.Close()
	infra.Logger.Sync()
	return nil
}
