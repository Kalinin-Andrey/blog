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

func (c *BlogReplicaSet) WriteRepo() blog.WriteFastRepository {
	return NewBlogRepository(c.ReplicaSet.WriteRepo())
}

func (c *BlogReplicaSet) ReadRepo() blog.ReadFastRepository {
	return NewBlogRepository(c.ReplicaSet.ReadRepo())
}
