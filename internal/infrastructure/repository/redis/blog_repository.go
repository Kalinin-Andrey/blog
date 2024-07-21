package redis

import (
	"github.com/minipkg/db/redis"
)

const (
	ratingPrefix = "blog_"
)

type BlogRepository struct {
	*Repository
	db redis.IDB
}

//var _ blog.WriteFastRepository = (*BlogRepository)(nil)
//var _ blog.ReadFastRepository = (*BlogRepository)(nil)
/*
func NewBlogRepository(repository *Repository) *BlogRepository {
	return &BlogRepository{
		Repository: repository,
	}
}

func (r *BlogRepository) ratingKey(sellerID uint) string {
	return ratingPrefix + strconv.Itoa(int(sellerID))
}

// Set создаём по SellerOldId в редисе запись с рейтингом
func (r *BlogRepository) Set(ctx context.Context, entity *blog.Blog) error {
	const metricName = "BlogRepository.Set"

	ratingB, err := blog_proto.Rating2RatingProto(entity).MarshalBinary()
	if err != nil {
		return fmt.Errorf("[%w] "+metricName+" MarshalBinary() error: %w", apperror.ErrInternal, err)
	}

	start := time.Now().UTC()
	err = r.DB().Set(ctx, r.ratingKey(entity.ID), string(ratingB), 0).Err()
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] "+metricName+" r.DB().Set() error: %w", apperror.ErrInternal, err)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *BlogRepository) MSet(ctx context.Context, ratingList *[]blog.Blog) error {
	const metricName = "BlogRepository.MSet"
	var item blog.Blog
	values := make([]interface{}, 0, len(*ratingList)*2)

	for _, item = range *ratingList {
		ratingB, err := blog_proto.Rating2RatingProto(&item).MarshalBinary()
		if err != nil {
			return fmt.Errorf("[%w] "+metricName+" MarshalBinary() error: %w", apperror.ErrInternal, err)
		}
		values = append(values, r.ratingKey(item.ID), string(ratingB))
	}

	start := time.Now().UTC()
	err := r.DB().MSet(ctx, values).Err()
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return fmt.Errorf("[%w] "+metricName+" r.DB().MSet() error: %w", apperror.ErrInternal, err)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

// Get получаем по SellerOldId из редиса запись с рейтингом
func (r *BlogRepository) Get(ctx context.Context, sellerID uint) (*blog.Blog, error) {
	const metricName = "BlogRepository.Get"
	start := time.Now().UTC()
	ratingProtoB, err := r.DB().Get(ctx, r.ratingKey(sellerID)).Bytes()
	if err != nil {
		if err.Error() == RedisNil {
			r.metrics.Inc(metricName, metricsSuccess)
			r.metrics.WriteTiming(start, metricName, metricsSuccess)
			return nil, apperror.ErrNotFound
		}
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] "+metricName+" r.DB().Get() error: %w", apperror.ErrInternal, err)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)

	ratingProto := &blog_proto.Blog{}
	err = ratingProto.UnmarshalBinary(ratingProtoB)
	if err != nil {
		return nil, fmt.Errorf("[%w] "+metricName+" ratingProto.UnmarshalBinary() error: %w", apperror.ErrInternal, err)
	}

	return blog_proto.RatingProto2Rating(ratingProto), nil
}

// MGetRating получаем по массиву SellerOldId из редиса записи с рейтингом
func (r *BlogRepository) MGet(ctx context.Context, sellerIDs *[]uint) (*[]blog.Blog, error) {
	const metricName = "BlogRepository.MGet"
	if sellerIDs == nil || len(*sellerIDs) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(*sellerIDs))
	for _, sellerID := range *sellerIDs {
		keys = append(keys, r.ratingKey(sellerID))
	}

	start := time.Now().UTC()
	res, err := r.DB().MGet(ctx, keys...).Result()
	if err != nil {
		if err.Error() == RedisNil {
			return nil, apperror.ErrNotFound
		}
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return nil, fmt.Errorf("[%w] "+metricName+" r.DB().Get() error: %w", apperror.ErrInternal, err)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)

	ratings := make([]blog.Blog, 0, len(*sellerIDs))
	for _, resItem := range res {
		ratingProtoS, ok := resItem.(string)
		if !ok {
			return nil, fmt.Errorf("[%w] "+metricName+" cast type resItem.(string) error", apperror.ErrInternal)
		}
		ratingProto := &blog_proto.Blog{}
		err = ratingProto.UnmarshalBinary([]byte(ratingProtoS))
		if err != nil {
			return nil, fmt.Errorf("[%w] "+metricName+" ratingProto.UnmarshalBinary() error: %w", apperror.ErrInternal, err)
		}

		r := blog_proto.RatingProto2Rating(ratingProto)
		ratings = append(ratings, *r)
	}

	return &ratings, nil
}
*/
