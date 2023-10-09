package blog

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RatingProto2Rating(ratingProto *Blog) (ratingDomain *blog.Blog) {
	if ratingProto == nil {
		return nil
	}

	ratingDomain = &blog.Blog{
		SellerOldId:         uint(ratingProto.SellerOldId),
		SellerName:          ratingProto.SellerName,
		Rating:              ratingProto.Rating,
		RatioDelivered:      ratingProto.RatioDelivered,
		RatioDefected:       ratingProto.RatioDefected,
		BuyerRatingWeight:   ratingProto.BuyerRatingWeight,
		AvgBuyerRating:      ratingProto.AvgBuyerRating,
		NbDelivered:         uint(ratingProto.NbDelivered),
		NbInDelivery:        uint(ratingProto.NbInDelivery),
		NbOrdersMarketplace: uint(ratingProto.NbOrdersMarketplace),
		NbBuyerRatings:      uint(ratingProto.NbBuyerRatings),
		NbDefected:          uint(ratingProto.NbDefected),
		NbOrdersTotal:       uint(ratingProto.NbOrdersTotal),
	}
	if ratingProto.Timestamp != nil && ratingProto.Timestamp.IsValid() {
		ratingDomain.Timestamp = ratingProto.Timestamp.AsTime()
	}
	return ratingDomain
}

func Rating2RatingProto(rating *blog.Blog) (ratingProto *Blog) {
	if rating == nil {
		return nil
	}

	return &Blog{
		Timestamp:           timestamppb.New(rating.Timestamp),
		SellerOldId:         uint64(rating.SellerOldId),
		SellerName:          rating.SellerName,
		Rating:              rating.Rating,
		RatioDelivered:      rating.RatioDelivered,
		RatioDefected:       rating.RatioDefected,
		BuyerRatingWeight:   rating.BuyerRatingWeight,
		AvgBuyerRating:      rating.AvgBuyerRating,
		NbDelivered:         uint64(rating.NbDelivered),
		NbInDelivery:        uint64(rating.NbInDelivery),
		NbOrdersMarketplace: uint64(rating.NbOrdersMarketplace),
		NbBuyerRatings:      uint64(rating.NbBuyerRatings),
		NbDefected:          uint64(rating.NbDefected),
		NbOrdersTotal:       uint64(rating.NbOrdersTotal),
	}
}
