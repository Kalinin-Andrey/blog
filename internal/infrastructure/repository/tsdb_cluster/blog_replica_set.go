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

func (c *BlogReplicaSet) WriteRepo() blog.WriteTsDBRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.WriteRepo())
}

func (c *BlogReplicaSet) ReadRepo() blog.ReadTsDBRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.ReadRepo())
}
