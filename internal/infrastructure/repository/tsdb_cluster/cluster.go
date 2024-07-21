package tsdb_cluster

import (
	"blog/internal/infrastructure/repository/tsdb"
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

func (c *Cluster) GetFirstShardMaster() *tsdb.Repository {
	return (*c.shards)[0].master
}

func (c *Cluster) GetShardMaster(n byte) *tsdb.Repository {
	return (*c.shards)[n].master
}

func (c *Cluster) GetShardSlave(n byte) *tsdb.Repository {
	return (*c.shards)[n].slave
}

func (c *Cluster) GetShardMasterByUintKey(sellerID uint) *tsdb.Repository {
	return (*c.GetShardByUintKey(sellerID)).WriteRepo()
}

func (c *Cluster) GetShardSlaveByUintKey(sellerID uint) *tsdb.Repository {
	return (*c.GetShardByUintKey(sellerID)).ReadRepo()
}

func (c *Cluster) GetShardMasterByStrKey(keyVal string) *tsdb.Repository {
	return (*c.GetShardByStrKey(keyVal)).WriteRepo()
}

func (c *Cluster) GetShardSlaveByStrKey(keyVal string) *tsdb.Repository {
	return (*c.GetShardByStrKey(keyVal)).ReadRepo()
}

func (c *Cluster) Close() {
	for i := range *c.shards {
		(*c.shards)[i].Close()
	}
}
