package redis

import (
	"context"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
	"github.com/pkg/errors"
	"github.com/wildberries-tech/wblogger"
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
	return &Cluster{
		shards:      shards,
		shardGetter: shardGetter,
	}
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

func (c *Cluster) GetFirstShardMaster() *Repository {
	return (*c.shards)[0].master
}

func (c *Cluster) GetShardMaster(n byte) *Repository {
	return (*c.shards)[n].master
}

func (c *Cluster) GetShardSlave(n byte) *Repository {
	return (*c.shards)[n].slave
}

func (c *Cluster) GetShardMasterByUintKey(keyVal uint) *Repository {
	return (*c.GetShardByUintKey(keyVal)).Master()
}

func (c *Cluster) GetShardSlaveByUintKey(keyVal uint) *Repository {
	return (*c.GetShardByUintKey(keyVal)).Slave()
}

func (c *Cluster) GetShardMasterByStrKey(keyVal string) *Repository {
	return (*c.GetShardByStrKey(keyVal)).Master()
}

func (c *Cluster) GetShardSlaveByStrKey(keyVal string) *Repository {
	return (*c.GetShardByStrKey(keyVal)).Slave()
}

func (c *Cluster) Close() {
	for i := range *c.shards {
		(*c.shards)[i].Close()
	}
}

func listenErrChanAndCancelIfError(ctx context.Context, metricName string, cancel context.CancelFunc, errChan chan error) {
	select {
	case <-ctx.Done():
		break
	case err := <-errChan:
		if err != nil {
			cancel()
			wblogger.Error(ctx, metricName+" error", errors.Wrap(apperror.ErrInternal, err.Error()))
			for err = range errChan {
				wblogger.Error(ctx, metricName+" error", errors.Wrap(apperror.ErrInternal, err.Error()))
			}
		}
		break
	}
	return
}

func (c *Cluster) Lock(ctx context.Context, key string) (bool, error) {
	return c.GetFirstShardMaster().Lock(ctx, key)
}
