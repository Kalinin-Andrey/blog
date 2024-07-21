package redis

type BlogCluster struct {
	*Cluster
}

/*
var _ blog.FastCluster = (*BlogCluster)(nil)

func NewBlogCluster(cluster *Cluster) *BlogCluster {
	return &BlogCluster{
		Cluster: cluster,
	}
}

func (c *BlogCluster) GetShardWriteRepo(num byte) blog.WriteFastRepository {
	return NewBlogRepository(c.Cluster.GetShardWriteRepo(num))
}

func (c *BlogCluster) GetShardReadRepo(num byte) blog.ReadFastRepository {
	return NewBlogRepository(c.Cluster.GetShardReadRepo(num))
}

func (c *BlogCluster) GetShardWriteRepoByID(ID uint) blog.WriteFastRepository {
	return NewBlogRepository(c.Cluster.GetShardWriteRepoByUintKey(ID))
}

func (c *BlogCluster) GetShardReadRepoByID(ID uint) blog.ReadFastRepository {
	return NewBlogRepository(c.Cluster.GetShardReadRepoByUintKey(ID))
}
*/
