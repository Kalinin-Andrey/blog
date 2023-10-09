package tsdb_cluster

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb"
)

type BlogCluster struct {
	*Cluster
}

var _ blog.TsDBCluster = (*BlogCluster)(nil)

func NewRatingCluster(cluster *Cluster) *BlogCluster {
	return &BlogCluster{
		Cluster: cluster,
	}
}

func (c *BlogCluster) GetShardMaster(num byte) blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardMaster(num))
}

func (c *BlogCluster) GetShardSlave(num byte) blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardSlave(num))
}

func (c *BlogCluster) GetShardMasterByID(ID uint) blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardMasterByUintKey(ID))
}

func (c *BlogCluster) GetShardSlaveByID(ID uint) blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardSlaveByUintKey(ID))
}
