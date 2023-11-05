package pg_cluster

import (
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/pg"
)

const (
	defaultSliceLen = 1000
)

type shardGetter interface {
	GetShardUint(key uint) byte
	GetShardStr(key string) byte
}

type Cluster struct {
	shards      *map[byte]*ReplicaSet
	shardGetter shardGetter
}

func NewCluster(shards *map[byte]*ReplicaSet, shardGetter shardGetter) *Cluster {
	return &Cluster{shards: shards, shardGetter: shardGetter}
}

func (c *Cluster) GetShardsNum() byte {
	return byte(len(*c.shards))
}

func (c *Cluster) GetShardByUintKey(key uint) *ReplicaSet {
	return (*c.shards)[c.shardGetter.GetShardUint(key)]
}

func (c *Cluster) GetShardByStrKey(key string) *ReplicaSet {
	return (*c.shards)[c.shardGetter.GetShardStr(key)]
}

func (c *Cluster) GetFirstShardWriteRepo() *pg.Repository {
	return (*c.shards)[0].master
}

func (c *Cluster) GetShardWriteRepo(n byte) *pg.Repository {
	return (*c.shards)[n].master
}

func (c *Cluster) GetShardReadRepo(n byte) *pg.Repository {
	return (*c.shards)[n].slave
}

func (c *Cluster) GetShardWriteRepoByUintKey(sellerID uint) *pg.Repository {
	return (*c.GetShardByUintKey(sellerID)).WriteRepo()
}

func (c *Cluster) GetShardReadRepoByUintKey(sellerID uint) *pg.Repository {
	return (*c.GetShardByUintKey(sellerID)).ReadRepo()
}

func (c *Cluster) GetShardWriteRepoByStrKey(keyVal string) *pg.Repository {
	return (*c.GetShardByStrKey(keyVal)).WriteRepo()
}

func (c *Cluster) GetShardReadRepoByStrKey(keyVal string) *pg.Repository {
	return (*c.GetShardByStrKey(keyVal)).ReadRepo()
}

func (c *Cluster) Close() {
	for i := range *c.shards {
		(*c.shards)[i].Close()
	}
}
