package tsdb_cluster

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/infrastructure/repository/tsdb"
)

type BlogReplicaSet struct {
	*ReplicaSet
}

var _ blog.TsDBReplicaSet = (*BlogReplicaSet)(nil)

func NewBlogReplicaSet(replicaSet *ReplicaSet) *BlogReplicaSet {
	return &BlogReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *BlogReplicaSet) Master() blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.Master())
}

func (c *BlogReplicaSet) Slave() blog.TsDBRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.Slave())
}
