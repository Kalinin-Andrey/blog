package tsdb_cluster

import (
	"blog/internal/domain/keyword"
	"blog/internal/infrastructure/repository/tsdb"
)

type KeywordReplicaSet struct {
	*ReplicaSet
}

var _ keyword.ReplicaSet = (*KeywordReplicaSet)(nil)

func NewKeywordReplicaSet(replicaSet *ReplicaSet) *KeywordReplicaSet {
	return &KeywordReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *KeywordReplicaSet) WriteRepo() keyword.WriteRepository {
	return tsdb.NewKeywordRepository(c.ReplicaSet.WriteRepo())
}

func (c *KeywordReplicaSet) ReadRepo() keyword.ReadRepository {
	return tsdb.NewKeywordRepository(c.ReplicaSet.ReadRepo())
}
