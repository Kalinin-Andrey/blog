package pg_cluster

import (
	"context"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/pg"
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

func (c *Cluster) GetFirstShardMaster() *pg.Repository {
	return (*c.shards)[0].master
}

func (c *Cluster) GetShardMaster(n byte) *pg.Repository {
	return (*c.shards)[n].master
}

func (c *Cluster) GetShardSlave(n byte) *pg.Repository {
	return (*c.shards)[n].slave
}

func (c *Cluster) GetShardMasterByUintKey(sellerID uint) *pg.Repository {
	return (*c.GetShardByUintKey(sellerID)).Master()
}

func (c *Cluster) GetShardSlaveByUintKey(sellerID uint) *pg.Repository {
	return (*c.GetShardByUintKey(sellerID)).Slave()
}

func (c *Cluster) GetShardMasterByStrKey(keyVal string) *pg.Repository {
	return (*c.GetShardByStrKey(keyVal)).Master()
}

func (c *Cluster) GetShardSlaveByStrKey(keyVal string) *pg.Repository {
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

func (c *Cluster) BeginOnSlaves(ctx context.Context) (pg.Txs, error) {
	var err error
	res := make(pg.Txs, int(c.GetShardsNum()))
	for n, replicaSet := range *c.shards {
		if res[n], err = replicaSet.slave.Begin(ctx); err != nil {
			return nil, errors.Wrap(apperror.ErrInternal, "slave.Begin() error: "+err.Error())
		}
	}
	return res, nil
}
