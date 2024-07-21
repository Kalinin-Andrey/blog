package infrastructure

import (
	"log"

	"github.com/minipkg/db/redis"

	"blog/internal/infrastructure/repository/pg"
	"blog/internal/infrastructure/repository/tsdb"
)

type Config struct {
	TsDB *TsDBConfig
	//Pg    *PgConfig
	Redis *RedisConfig
}

type TsDBConfig struct {
	Conn       *tsdb.Config
	ReplicaSet *pgReplicaSet
}

type PgConfig struct {
	Conn   *pg.Config
	Shards *map[byte]pgReplicaSet
}

type pgReplicaSet struct {
	MasterHost string
	MasterPort string
	SlaveHost  string
	SlavePort  string
	Password   string
}

// ключ - номер шарда, 0 индекс- мастер, 1 индекс - слейв
func (c *PgConfig) getConfig() map[byte][2]pg.Config {
	res := make(map[byte][2]pg.Config)

	if c.Conn.User == "" {
		log.Fatal("Empty Infra.Pg.Conn.User")
	}

	for k, v := range *c.Shards {
		password := v.Password
		if password == "" {
			log.Fatalf("Empty Infra.Pg.Shards.%d.Password", k)
		}
		shardConf := [2]pg.Config{}
		shardConf[0] = pg.Config{
			Host:            v.MasterHost,
			Port:            v.MasterPort,
			User:            c.Conn.User,
			Password:        password,
			DbName:          c.Conn.DbName,
			SchemaName:      c.Conn.SchemaName,
			MaxOpenConns:    c.Conn.MaxOpenConns,
			MaxIdleConns:    c.Conn.MaxIdleConns,
			MaxConnLifetime: c.Conn.MaxConnLifetime,
			Timeout:         c.Conn.Timeout,
		}

		shardConf[1] = pg.Config{
			Host:            v.SlaveHost,
			Port:            v.SlavePort,
			User:            c.Conn.User,
			Password:        password,
			DbName:          c.Conn.DbName,
			SchemaName:      c.Conn.SchemaName,
			MaxOpenConns:    c.Conn.MaxOpenConns,
			MaxIdleConns:    c.Conn.MaxIdleConns,
			MaxConnLifetime: c.Conn.MaxConnLifetime,
			Timeout:         c.Conn.Timeout,
		}

		res[k] = shardConf
	}

	return res
}

func (c *TsDBConfig) getConfig() *[2]tsdb.Config {
	var res [2]tsdb.Config

	if c.Conn.User == "" {
		log.Fatal("Empty Infra.TsDB.Conn.User")
	}

	res[0] = tsdb.Config{
		Host:            c.ReplicaSet.MasterHost,
		Port:            c.ReplicaSet.MasterPort,
		User:            c.Conn.User,
		Password:        c.ReplicaSet.Password,
		DbName:          c.Conn.DbName,
		SchemaName:      c.Conn.SchemaName,
		MaxOpenConns:    c.Conn.MaxOpenConns,
		MaxIdleConns:    c.Conn.MaxIdleConns,
		MaxConnLifetime: c.Conn.MaxConnLifetime,
		Timeout:         c.Conn.Timeout,
	}

	res[1] = tsdb.Config{
		Host:            c.ReplicaSet.SlaveHost,
		Port:            c.ReplicaSet.SlavePort,
		User:            c.Conn.User,
		Password:        c.ReplicaSet.Password,
		DbName:          c.Conn.DbName,
		SchemaName:      c.Conn.SchemaName,
		MaxOpenConns:    c.Conn.MaxOpenConns,
		MaxIdleConns:    c.Conn.MaxIdleConns,
		MaxConnLifetime: c.Conn.MaxConnLifetime,
		Timeout:         c.Conn.Timeout,
	}

	return &res
}

type RedisConfig struct {
	Conn       *redis.Config
	ReplicaSet *redisReplicaSet
}

type redisReplicaSet struct {
	MasterAddrs string
	SlaveAddrs  string
	Password    string
}

// ключ - номер шарда, 0 индекс- мастер, 1 индекс - слейв
func (c *RedisConfig) getConfig() *[2]redis.Config {
	var res [2]redis.Config

	res[0] = redis.Config{
		Addrs:    []string{c.ReplicaSet.MasterAddrs},
		Login:    c.Conn.Login,
		Password: c.ReplicaSet.Password,
		DBNum:    c.Conn.DBNum,
	}

	res[1] = redis.Config{
		Addrs:    []string{c.ReplicaSet.SlaveAddrs},
		Login:    c.Conn.Login,
		Password: c.ReplicaSet.Password,
		DBNum:    c.Conn.DBNum,
	}

	return &res
}
