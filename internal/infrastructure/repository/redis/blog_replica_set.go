package redis

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
)

type BlogReplicaSet struct {
	*ReplicaSet
}

var _ blog.FastReplicaSet = (*BlogReplicaSet)(nil)

func NewBlogReplicaSet(replicaSet *ReplicaSet) *BlogReplicaSet {
	return &BlogReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *BlogReplicaSet) Master() blog.FastRepository {
	return NewBlogRepository(c.ReplicaSet.Master())
}

func (c *BlogReplicaSet) Slave() blog.FastRepository {
	return NewBlogRepository(c.ReplicaSet.Slave())
}
