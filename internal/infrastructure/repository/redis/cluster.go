package redis

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
	return &Cluster{
		shards:      shards,
		shardGetter: shardGetter,
	}
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

func (c *Cluster) GetFirstShardWriteRepo() *Repository {
	return (*c.shards)[0].master
}

func (c *Cluster) GetShardWriteRepo(n byte) *Repository {
	return (*c.shards)[n].master
}

func (c *Cluster) GetShardReadRepo(n byte) *Repository {
	return (*c.shards)[n].slave
}

func (c *Cluster) GetShardWriteRepoByUintKey(keyVal uint) *Repository {
	return (*c.GetShardByUintKey(keyVal)).WriteRepo()
}

func (c *Cluster) GetShardReadRepoByUintKey(keyVal uint) *Repository {
	return (*c.GetShardByUintKey(keyVal)).ReadRepo()
}

func (c *Cluster) GetShardWriteRepoByStrKey(keyVal string) *Repository {
	return (*c.GetShardByStrKey(keyVal)).WriteRepo()
}

func (c *Cluster) GetShardReadRepoByStrKey(keyVal string) *Repository {
	return (*c.GetShardByStrKey(keyVal)).ReadRepo()
}

func (c *Cluster) Close() {
	for i := range *c.shards {
		(*c.shards)[i].Close()
	}
}
