package tsdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Kalinin-Andrey/blog/internal/domain"
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
	"github.com/pkg/errors"
	"github.com/wildberries-tech/wblogger"
	"strconv"
	"strings"
	"time"
)

type BlogRepository struct {
	*Repository
}

var _ blog.WriteTsDBRepository = (*BlogRepository)(nil)
var _ blog.ReadTsDBRepository = (*BlogRepository)(nil)

func NewBlogRepository(repository *Repository) *BlogRepository {
	return &BlogRepository{
		Repository: repository,
	}
}

const (
	rating_field_SellerOldID    = "seller_old_id"
	rating_field_Rating         = "rating"
	rating_field_AvgBuyerRating = "avg_buyer_rating"
	rating_field_RatioDefected  = "ratio_defected"

	rating_sql_GetLast                    = "SELECT seller_old_id, rating, ratio_delivered, nb_delivered, nb_in_delivery, nb_orders_marketplace, buyer_rating_weight, avg_buyer_rating, nb_buyer_ratings, ratio_defected, nb_defected, nb_orders_total, timestamp FROM rating.rating WHERE seller_id = $1;"
	rating_sql_Create                     = "INSERT INTO rating.rating(seller_old_id, rating, ratio_delivered, nb_delivered, nb_in_delivery, nb_orders_marketplace, buyer_rating_weight, avg_buyer_rating, nb_buyer_ratings, ratio_defected, nb_defected, nb_orders_total, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) "
	rating_sql_Create_OnConflictDoUpdate  = " ON CONFLICT (seller_old_id) DO UPDATE SET rating = EXCLUDED.rating, ratio_delivered = EXCLUDED.ratio_delivered, nb_delivered = EXCLUDED.nb_delivered, nb_in_delivery = EXCLUDED.nb_in_delivery, nb_orders_marketplace = EXCLUDED.nb_orders_marketplace, buyer_rating_weight = EXCLUDED.buyer_rating_weight, avg_buyer_rating = EXCLUDED.avg_buyer_rating, nb_buyer_ratings = EXCLUDED.nb_buyer_ratings, ratio_defected = EXCLUDED.ratio_defected, nb_defected = EXCLUDED.nb_defected, nb_orders_total = EXCLUDED.nb_orders_total, timestamp = EXCLUDED.timestamp;"
	rating_sql_MCreate                    = "INSERT INTO rating.rating (seller_old_id, rating, ratio_delivered, nb_delivered, nb_in_delivery, nb_orders_marketplace, buyer_rating_weight, avg_buyer_rating, nb_buyer_ratings, ratio_defected, nb_defected, nb_orders_total, timestamp) VALUES "
	rating_sql_MGet                       = "SELECT seller_old_id, rating, ratio_delivered, nb_delivered, nb_in_delivery, nb_orders_marketplace, buyer_rating_weight, avg_buyer_rating, nb_buyer_ratings, ratio_defected, nb_defected, nb_orders_total, timestamp FROM rating.rating ${where} ORDER BY ${orderBy} LIMIT $${limit};"
	rating_sql_SellerOldID_AnyCondition   = " seller_old_id = ANY($${i}) "
	rating_sql_SellerOldID_LCondition     = " seller_old_id < $${i} "
	rating_sql_SellerOldID_GCondition     = " seller_old_id > $${i} "
	rating_sql_Rating_LECondition         = " (rating < $${i} OR (rating = $${i} ${SellerOldID_GCondition})) "
	rating_sql_Rating_GECondition         = " (rating > $${i} OR (rating = $${i} ${SellerOldID_GCondition})) "
	rating_sql_AvgBuyerRating_LECondition = " (avg_buyer_rating < $${i} OR (avg_buyer_rating = $${i} ${SellerOldID_GCondition})) "
	rating_sql_AvgBuyerRating_GECondition = " (avg_buyer_rating > $${i} OR (avg_buyer_rating = $${i} ${SellerOldID_GCondition})) "
	rating_sql_RatioDefected_LECondition  = " (ratio_delivered < $${i}  OR (ratio_delivered = $${i} ${SellerOldID_GCondition})) "
	rating_sql_RatioDefected_GECondition  = " (ratio_delivered > $${i} OR (ratio_delivered = $${i} ${SellerOldID_GCondition})) "
)

func (r *BlogRepository) Begin(ctx context.Context) (domain.Tx, error) {
	return r.Repository.Begin(ctx)
}

func (r *BlogRepository) GetLast(ctx context.Context, sellerID uint) (*blog.Blog, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	const metricName = "BlogRepository.GetLast"
	start := time.Now().UTC()

	entity := &blog.Blog{}
	if err := r.db.QueryRow(ctx, rating_sql_GetLast, sellerID).Scan(&entity.SellerOldId, &entity.Rating, &entity.RatioDelivered, &entity.NbDelivered, &entity.NbInDelivery, &entity.NbOrdersMarketplace, &entity.BuyerRatingWeight, &entity.AvgBuyerRating, &entity.NbBuyerRatings, &entity.RatioDefected, &entity.NbDefected, &entity.NbOrdersTotal, &entity.Timestamp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		wblogger.Error(ctx, metricName+"-Query", err)
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return entity, nil
}

func (r *BlogRepository) Create(ctx context.Context, entity *blog.Blog) error {
	const metricName = "BlogRepository.Create"
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, rating_sql_Create+rating_sql_Create_OnConflictDoUpdate, entity.SellerOldId, entity.Rating, entity.RatioDelivered, entity.NbDelivered, entity.NbInDelivery, entity.NbOrdersMarketplace, entity.BuyerRatingWeight, entity.AvgBuyerRating, entity.NbBuyerRatings, entity.RatioDefected, entity.NbDefected, entity.NbOrdersTotal, entity.Timestamp)
	if err != nil {
		wblogger.Error(ctx, metricName+" Exec: "+rating_sql_Create+rating_sql_Create_OnConflictDoUpdate, err)
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return errors.Wrap(apperror.ErrInternal, err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *BlogRepository) MCreate(ctx context.Context, entities *[]blog.Blog) error {
	const metricName = "BlogRepository.MCreate"
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	if len(*entities) == 0 {
		return nil
	}
	b := strings.Builder{}
	params := make([]interface{}, 0, len(*entities)*13)
	b.WriteString(rating_sql_MCreate)
	for i, entity := range *entities {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString("($" + strconv.Itoa(i*13+1) + ", $" + strconv.Itoa(i*13+2) + ", $" + strconv.Itoa(i*13+3) + ", $" + strconv.Itoa(i*13+4) + ", $" + strconv.Itoa(i*13+5) + ", $" + strconv.Itoa(i*13+6) + ", $" + strconv.Itoa(i*13+7) + ", $" + strconv.Itoa(i*13+8) + ", $" + strconv.Itoa(i*13+9) + ", $" + strconv.Itoa(i*13+10) + ", $" + strconv.Itoa(i*13+11) + ", $" + strconv.Itoa(i*13+12) + ", $" + strconv.Itoa(i*13+13) + ")")
		params = append(params, entity.SellerOldId, entity.Rating, entity.RatioDelivered, entity.NbDelivered, entity.NbInDelivery, entity.NbOrdersMarketplace, entity.BuyerRatingWeight, entity.AvgBuyerRating, entity.NbBuyerRatings, entity.RatioDefected, entity.NbDefected, entity.NbOrdersTotal, entity.Timestamp)
	}
	b.WriteString(rating_sql_Create_OnConflictDoUpdate)
	start := time.Now().UTC()

	_, err := r.db.Exec(ctx, b.String(), params...)
	if err != nil {
		wblogger.Error(ctx, metricName+"-Exec sql: "+b.String()+" error: ", err)
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return errors.Wrap(apperror.ErrInternal, err.Error())
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)
	return nil
}

func (r *BlogRepository) sqlQueryForMGet(filter *blog.Filter4TsDB) (string, *[]interface{}) {
	sql := rating_sql_MGet
	var whereBuilder strings.Builder
	var orderByBuilder strings.Builder
	var i int
	params := make([]interface{}, 0, 4)
	var isSortBySellerOldID bool
	var isSortByRating bool
	var isSortByAvgBuyerRating bool
	var isSortByRatioDefected bool

	if filter.SortBy != nil {
		switch *filter.SortBy {
		case blog.JsonProperty_SellerOldID:
			isSortBySellerOldID = true
		case blog.JsonProperty_Rating:
			isSortByRating = true
		case blog.JsonProperty_AvgBuyerRating:
			isSortByAvgBuyerRating = true
		case blog.JsonProperty_RatioDefected:
			isSortByRatioDefected = true
		}
	}

	if filter.SellerOldIDs != nil {
		i++
		params = append(params, *filter.SellerOldIDs)
		whereBuilder.WriteString(strings.ReplaceAll(rating_sql_SellerOldID_AnyCondition, "${i}", strconv.Itoa(i)))
	}

	if isSortBySellerOldID && filter.FromSellerOldID != nil {
		i++
		params = append(params, *filter.FromSellerOldID)

		if whereBuilder.Len() > 0 {
			whereBuilder.WriteString(sql_And)
		}

		if filter.IsSortOrderDesc {
			whereBuilder.WriteString(strings.ReplaceAll(rating_sql_SellerOldID_LCondition, "${i}", strconv.Itoa(i)))
		} else {
			whereBuilder.WriteString(strings.ReplaceAll(rating_sql_SellerOldID_GCondition, "${i}", strconv.Itoa(i)))
		}
	}

	if filter.FromRating != nil {
		i++
		params = append(params, *filter.FromRating)

		if whereBuilder.Len() > 0 {
			whereBuilder.WriteString(sql_And)
		}

		s := ""
		if isSortByRating && filter.IsSortOrderDesc {
			s = strings.ReplaceAll(rating_sql_Rating_LECondition, "${i}", strconv.Itoa(i))
		} else {
			s = strings.ReplaceAll(rating_sql_Rating_GECondition, "${i}", strconv.Itoa(i))
		}
		if filter.FromSellerOldID != nil {
			i++
			params = append(params, *filter.FromSellerOldID)
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", sql_And+strings.ReplaceAll(rating_sql_SellerOldID_GCondition, "${i}", strconv.Itoa(i)))
		} else {
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", "")
		}
		whereBuilder.WriteString(s)
	}

	if filter.FromAvgBuyerRating != nil {
		i++
		params = append(params, *filter.FromAvgBuyerRating)

		if whereBuilder.Len() > 0 {
			whereBuilder.WriteString(sql_And)
		}

		s := ""
		if isSortByAvgBuyerRating && filter.IsSortOrderDesc {
			s = strings.ReplaceAll(rating_sql_AvgBuyerRating_LECondition, "${i}", strconv.Itoa(i))
		} else {
			s = strings.ReplaceAll(rating_sql_AvgBuyerRating_GECondition, "${i}", strconv.Itoa(i))
		}
		if filter.FromSellerOldID != nil {
			i++
			params = append(params, *filter.FromSellerOldID)
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", sql_And+strings.ReplaceAll(rating_sql_SellerOldID_GCondition, "${i}", strconv.Itoa(i)))
		} else {
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", "")
		}
		whereBuilder.WriteString(s)
	}

	if filter.FromRatioDefected != nil {
		i++
		params = append(params, *filter.FromRatioDefected)

		if whereBuilder.Len() > 0 {
			whereBuilder.WriteString(sql_And)
		}

		s := ""
		if isSortByRatioDefected && filter.IsSortOrderDesc {
			s = strings.ReplaceAll(rating_sql_RatioDefected_LECondition, "${i}", strconv.Itoa(i))
		} else {
			s = strings.ReplaceAll(rating_sql_RatioDefected_GECondition, "${i}", strconv.Itoa(i))
		}
		if filter.FromSellerOldID != nil {
			i++
			params = append(params, *filter.FromSellerOldID)
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", sql_And+strings.ReplaceAll(rating_sql_SellerOldID_GCondition, "${i}", strconv.Itoa(i)))
		} else {
			s = strings.ReplaceAll(s, "${SellerOldID_GCondition}", "")
		}
		whereBuilder.WriteString(s)
	}

	if isSortBySellerOldID {
		orderByBuilder.WriteString(rating_field_SellerOldID)

		if filter.IsSortOrderDesc {
			orderByBuilder.WriteString(sql_Desc)
		}
	}

	if isSortByRating {
		if orderByBuilder.Len() > 0 {
			orderByBuilder.WriteString(", ")
		}
		orderByBuilder.WriteString(rating_field_Rating)

		if filter.IsSortOrderDesc {
			orderByBuilder.WriteString(sql_Desc)
		}
		orderByBuilder.WriteString(", " + rating_field_SellerOldID)
	}

	if isSortByAvgBuyerRating {
		if orderByBuilder.Len() > 0 {
			orderByBuilder.WriteString(", ")
		}
		orderByBuilder.WriteString(rating_field_AvgBuyerRating)

		if filter.IsSortOrderDesc {
			orderByBuilder.WriteString(sql_Desc)
		}
		orderByBuilder.WriteString(", " + rating_field_SellerOldID)
	}

	if isSortByRatioDefected {
		if orderByBuilder.Len() > 0 {
			orderByBuilder.WriteString(", ")
		}
		orderByBuilder.WriteString(rating_field_RatioDefected)

		if filter.IsSortOrderDesc {
			orderByBuilder.WriteString(sql_Desc)
		}
		orderByBuilder.WriteString(", " + rating_field_SellerOldID)
	}

	if whereBuilder.Len() > 0 {
		sql = strings.ReplaceAll(sql, "${where}", sql_Where+whereBuilder.String())
	} else {
		sql = strings.ReplaceAll(sql, "${where}", "")
	}

	sql = strings.ReplaceAll(sql, "${orderBy}", orderByBuilder.String())
	i++
	params = append(params, *filter.Limit)
	sql = strings.ReplaceAll(sql, "${limit}", strconv.Itoa(i))

	return sql, &params
}

func (r *BlogRepository) MGet(ctx context.Context, filter *blog.Filter4TsDB) (*[]blog.Blog, error) {
	//ctx, cancel := context.WithTimeout(ctx, r.timeout)
	//defer cancel()
	var rows pgx.Rows
	var err error
	const metricName = "BlogRepository.MGet"

	limit := uint(blog.Default_Limit)
	if filter.Limit != nil {
		limit = *filter.Limit
	}
	if filter.SellerOldIDs != nil && uint(len(*filter.SellerOldIDs)) < limit {
		limit = uint(len(*filter.SellerOldIDs))
	}

	sql, params := r.sqlQueryForMGet(filter)

	start := time.Now().UTC()
	rows, err = r.db.Query(ctx, sql, *params...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		wblogger.Error(ctx, metricName+"-Query", err)
		r.metrics.Inc(metricName, metricsFail)
		r.metrics.WriteTiming(start, metricName, metricsFail)
		return nil, errors.Wrap(apperror.ErrInternal, err.Error())
	}
	defer rows.Close()

	res := make([]blog.Blog, 0, limit)
	entity := blog.Blog{}

	for rows.Next() {
		if err = rows.Scan(&entity.SellerOldId, &entity.Rating, &entity.RatioDelivered, &entity.NbDelivered, &entity.NbInDelivery, &entity.NbOrdersMarketplace, &entity.BuyerRatingWeight, &entity.AvgBuyerRating, &entity.NbBuyerRatings, &entity.RatioDefected, &entity.NbDefected, &entity.NbOrdersTotal, &entity.Timestamp); err != nil {
			wblogger.Error(ctx, metricName+"-rows.Scan", err)
			r.metrics.Inc(metricName, metricsFail)
			r.metrics.WriteTiming(start, metricName, metricsFail)
			return nil, errors.Wrap(apperror.ErrInternal, err.Error())
		}
		res = append(res, entity)
	}
	r.metrics.Inc(metricName, metricsSuccess)
	r.metrics.WriteTiming(start, metricName, metricsSuccess)

	if len(res) == 0 { // да, и такое почему-то случается
		return nil, apperror.ErrNotFound
	}
	return &res, nil
}

func getTableFieldByPropJson(prop string) (string, error) {
	var res string
	switch prop {
	case "sellerOldId":
		res = rating_field_SellerOldID
	case "rating":
		res = rating_field_Rating
	case "avgBuyerRating":
		res = rating_field_AvgBuyerRating
	case "ratioDefected":
		res = rating_field_RatioDefected
	default:
		return "", fmt.Errorf("[%w] Wrong param sortBy", apperror.ErrBadRequest)
	}
	return res, nil
}
