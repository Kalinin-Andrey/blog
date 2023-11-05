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

func (c *BlogCluster) GetShardWriteRepo(num byte) blog.WriteTsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardMaster(num))
}

func (c *BlogCluster) GetShardReadRepo(num byte) blog.ReadTsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardSlave(num))
}

func (c *BlogCluster) GetShardWriteRepoByID(ID uint) blog.WriteTsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardMasterByUintKey(ID))
}

func (c *BlogCluster) GetShardReadRepoByID(ID uint) blog.ReadTsDBRepository {
	return tsdb.NewBlogRepository(c.Cluster.GetShardSlaveByUintKey(ID))
}
