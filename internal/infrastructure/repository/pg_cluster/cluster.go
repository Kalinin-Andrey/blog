package pg_cluster

import (
	"blog/internal/infrastructure/repository/pg"
)

const (
	defaultCapacityForResult = 100
)

type shardGetter interface {
	GetShardUint(key uint) byte
	GetShardStr(key string) byte
}

type Cluster struct {
	shards      map[byte]*ReplicaSet
	shardGetter shardGetter
}

func NewCluster(shards map[byte]*ReplicaSet, shardGetter shardGetter) *Cluster {
	return &Cluster{shards: shards, shardGetter: shardGetter}
}

func (c *Cluster) GetShardsNum() byte {
	return byte(len(c.shards))
}

func (c *Cluster) GetShardByUint(key uint) *ReplicaSet {
	return c.shards[c.shardGetter.GetShardUint(key)]
}

func (c *Cluster) GetShardByStrKey(key string) *ReplicaSet {
	return c.shards[c.shardGetter.GetShardStr(key)]
}

func (c *Cluster) WriteRepoFirstShard() *pg.Repository {
	return c.shards[1].master
}

func (c *Cluster) WriteRepoByNum(n byte) *pg.Repository {
	return c.shards[n].master
}

func (c *Cluster) ReadRepoByNum(n byte) *pg.Repository {
	return c.shards[n].slave
}

func (c *Cluster) WriteRepoByUint(ID uint) *pg.Repository {
	return c.GetShardByUint(ID).WriteRepo()
}

func (c *Cluster) ReadRepoByUint(ID uint) *pg.Repository {
	return c.GetShardByUint(ID).ReadRepo()
}

func (c *Cluster) WriteRepoByStr(keyVal string) *pg.Repository {
	return (*c.GetShardByStrKey(keyVal)).WriteRepo()
}

func (c *Cluster) ReadRepoByStr(keyVal string) *pg.Repository {
	return (*c.GetShardByStrKey(keyVal)).ReadRepo()
}

func (c *Cluster) Close() {
	for i := range c.shards {
		c.shards[i].Close()
	}
}
