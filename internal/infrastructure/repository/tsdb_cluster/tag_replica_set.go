package tsdb_cluster

import (
	"blog/internal/domain/tag"
	"blog/internal/infrastructure/repository/tsdb"
)

type TagReplicaSet struct {
	*ReplicaSet
}

var _ tag.ReplicaSet = (*TagReplicaSet)(nil)

func NewTagReplicaSet(replicaSet *ReplicaSet) *TagReplicaSet {
	return &TagReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *TagReplicaSet) WriteRepo() tag.WriteRepository {
	return tsdb.NewTagRepository(c.ReplicaSet.WriteRepo())
}

func (c *TagReplicaSet) ReadRepo() tag.ReadRepository {
	return tsdb.NewTagRepository(c.ReplicaSet.ReadRepo())
}
