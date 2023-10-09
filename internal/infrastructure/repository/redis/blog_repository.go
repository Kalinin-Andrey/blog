package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/minipkg/db/redis"
	"github.com/pkg/errors"

	blog_proto "github.com/Kalinin-Andrey/blog/internal/app/proto/blog"
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
)

const (
	ratingPrefix = "blog_"
)

type BlogRepository struct {
	*Repository
	db redis.IDB
}

var _ blog.FastRepository = (*BlogRepository)(nil)

func NewBlogRepository(repository *Repository) *BlogRepository {
	return &BlogRepository{
		Repository: repository,
	}
}

func (r *BlogRepository) ratingKey(sellerID uint) string {
	return ratingPrefix + strconv.Itoa(int(sellerID))
}

// Set создаём по SellerOldId в редисе запись с рейтингом
func (r *BlogRepository) Set(ctx context.Context, rating *blog.Blog) error {
	const metricName = "BlogRepository.Set"

	ratingB, err := blog_proto.Rating2RatingProto(rating).MarshalBinary()
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, metricName+" MarshalBinary() error: %s", err)
	}

	start := time.Now().UTC()
	err = r.DB().Set(ctx, r.ratingKey(rating.SellerOldId), string(ratingB), 0).Err()
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return errors.Wrapf(apperror.ErrInternal, metricName+" r.DB().Set() error: %s", err)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *BlogRepository) MSet(ctx context.Context, ratingList *[]blog.Blog) error {
	const metricName = "BlogRepository.MSet"
	var ratingItem blog.Blog
	values := make([]interface{}, 0, len(*ratingList)*2)

	for _, ratingItem = range *ratingList {
		ratingB, err := blog_proto.Rating2RatingProto(&ratingItem).MarshalBinary()
		if err != nil {
			return errors.Wrapf(apperror.ErrInternal, metricName+" MarshalBinary() error: %s", err)
		}
		values = append(values, r.ratingKey(ratingItem.SellerOldId), string(ratingB))
	}

	start := time.Now().UTC()
	err := r.DB().MSet(ctx, values).Err()
	if err != nil {
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return errors.Wrapf(apperror.ErrInternal, metricName+" r.DB().MSet() error: %s", err)
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
		return nil, errors.Wrapf(apperror.ErrInternal, metricName+" r.DB().Get() error: %s", err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)

	ratingProto := &blog_proto.Blog{}
	err = ratingProto.UnmarshalBinary(ratingProtoB)
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrInternal, metricName+" ratingProto.UnmarshalBinary() error: %s", err.Error())
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
		return nil, errors.Wrapf(apperror.ErrInternal, metricName+" r.DB().Get() error: %s", err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)

	ratings := make([]blog.Blog, 0, len(*sellerIDs))
	for _, resItem := range res {
		ratingProtoS, ok := resItem.(string)
		if !ok {
			return nil, errors.Wrapf(apperror.ErrInternal, metricName+" cast type resItem.(string) error")
		}
		ratingProto := &blog_proto.Blog{}
		err = ratingProto.UnmarshalBinary([]byte(ratingProtoS))
		if err != nil {
			return nil, errors.Wrapf(apperror.ErrInternal, metricName+" ratingProto.UnmarshalBinary() error: %s", err.Error())
		}

		r := blog_proto.RatingProto2Rating(ratingProto)
		ratings = append(ratings, *r)
	}

	return &ratings, nil
}
