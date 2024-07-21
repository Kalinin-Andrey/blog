package infrastructure

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minipkg/prometheus-utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"blog/internal/infrastructure/repository/redis"
	"blog/internal/infrastructure/repository/tsdb"
	"blog/internal/infrastructure/repository/tsdb_cluster"
	_ "blog/migrations/tsdb"
	tsdbmigration "blog/migrations/tsdb"

	"blog/internal/pkg/config"
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
	case config.Environment_Local:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case config.Environment_Dev:
		lvl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case config.Environment_Stage:
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case config.Environment_Prd:
		//lvl = zap.NewAtomicLevelAt(zap.ErrorLevel)
		lvl = zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	infra.Logger, err = zap.Config{
		Level:             lvl,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			EncodeTime: zapcore.RFC3339TimeEncoder,
		},
		InitialFields: nil,
	}.Build()
	return err
}

func (infra *Infrastructure) tsDBInit(ctx context.Context) error {
	conf := infra.config.TsDB.getConfig()
	nodeLabel := "tsdb_master"
	sqlMetrics := prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	dbMetrics := prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	gaugeMetrics := prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[0].DbName)
	master, err := tsdb.NewRepository(conf[0], dbMetrics, sqlMetrics, gaugeMetrics)
	if err != nil {
		return err
	}

	nodeLabel = "tsdb_replicas"
	sqlMetrics = prometheus_utils.NewSqlMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	dbMetrics = prometheus_utils.NewDbMetrics(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	gaugeMetrics = prometheus_utils.NewDBGauge(infra.appConfig.NameSpace, infra.appConfig.Name, infra.appConfig.Service, nodeLabel, conf[1].DbName)
	slave, err := tsdb.NewRepository(conf[1], dbMetrics, sqlMetrics, gaugeMetrics)
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
