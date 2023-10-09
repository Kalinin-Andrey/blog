package redis

import (
	"context"
	"sync"

	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
)

type BlogCluster struct {
	*Cluster
}

var _ blog.FastCluster = (*BlogCluster)(nil)

func NewBlogCluster(cluster *Cluster) *BlogCluster {
	return &BlogCluster{
		Cluster: cluster,
	}
}

func (c *BlogCluster) GetShardMaster(num byte) blog.FastRepository {
	return NewBlogRepository(c.Cluster.GetShardMaster(num))
}

func (c *BlogCluster) GetShardSlave(num byte) blog.FastRepository {
	return NewBlogRepository(c.Cluster.GetShardSlave(num))
}

func (c *BlogCluster) GetShardMasterByID(ID uint) blog.FastRepository {
	return NewBlogRepository(c.Cluster.GetShardMasterByUintKey(ID))
}

func (c *BlogCluster) GetShardSlaveByID(ID uint) blog.FastRepository {
	return NewBlogRepository(c.Cluster.GetShardSlaveByUintKey(ID))
}

func (c *BlogCluster) MCreate(ctx context.Context, entities *[]blog.Blog) error {
	const metricName = "BlogCluster.MCreate"
	var entity blog.Blog
	shardsNum := c.GetShardsNum()
	entitiesByShardsNums := make(map[byte][]blog.Blog, shardsNum)

	for _, entity = range *entities {
		n := c.shardGetter.GetShardUint(entity.SellerOldId)
		if entitiesByShardsNums[n] == nil {
			entitiesByShardsNums[n] = make([]blog.Blog, 0, defaultSliceLen)
		}
		entitiesByShardsNums[n] = append(entitiesByShardsNums[n], entity)
	}

	errCh := make(chan error, shardsNum)
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go listenErrChanAndCancelIfError(ctx, metricName, cancel, errCh)
	for i := range entitiesByShardsNums {
		wg.Add(1)
		entitiesItem := entitiesByShardsNums[i]
		go c.mCreateItem(ctx, wg, i, &entitiesItem, errCh)
	}

	wg.Wait()
	close(errCh)

	return nil
}

func (c *BlogCluster) mCreateItem(ctx context.Context, wg *sync.WaitGroup, shardNum byte, entities *[]blog.Blog, errCh chan error) {
	defer wg.Done()
	err := c.GetShardMaster(shardNum).MSet(ctx, entities)
	if err != nil {
		errCh <- err
	}
	return
}

func (c *BlogCluster) MGet(ctx context.Context, sellerIDs *[]uint) (*[]blog.Blog, error) {
	const metricName = "BlogCluster.MGet"
	var id uint
	shardsNum := c.GetShardsNum()
	idsByShardsNums := make(map[byte][]uint, shardsNum)
	res := make([]blog.Blog, 0, len(*sellerIDs))

	for _, id = range *sellerIDs {
		n := c.shardGetter.GetShardUint(id)
		idsByShardsNums[n] = append(idsByShardsNums[n], id)
	}

	errCh := make(chan error, shardsNum)
	resCh := make(chan *[]blog.Blog, shardsNum)
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go listenErrChanAndCancelIfError(ctx, metricName, cancel, errCh)
	for n := range idsByShardsNums {
		wg.Add(1)
		ids := idsByShardsNums[n]
		go c.mGetItem(ctx, wg, n, &ids, resCh, errCh)
	}

	wg.Wait()
	close(errCh)
	close(resCh)

	for slp := range resCh {
		res = append(res, (*slp)...)
	}

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}

	return &res, nil
}

func (c *BlogCluster) mGetItem(ctx context.Context, wg *sync.WaitGroup, shardNum byte, sellerIDs *[]uint, resCh chan *[]blog.Blog, errCh chan error) {
	defer wg.Done()
	res, err := c.GetShardMaster(shardNum).MGet(ctx, sellerIDs)
	if err != nil {
		errCh <- err
		return
	}
	resCh <- res
	return
}
