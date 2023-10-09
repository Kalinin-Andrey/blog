package controller

import (
	"github.com/Kalinin-Andrey/blog/internal/domain/blog"
	"github.com/Kalinin-Andrey/blog/internal/pkg/fasthttp_tools"
	routing "github.com/qiangxue/fasthttp-routing"
)

type blogController struct {
	router  *routing.Router
	service *blog.Service
}

func NewBlogController(router *routing.Router, service *blog.Service) *blogController {
	return &blogController{
		router:  router,
		service: service,
	}
}

/*
func (c *blogController) Get(rctx *routing.Context) error {
	const metricName = "blogController.Get"
	var res *fasthttp_tools.Response

	sellerID, err := strconv.ParseUint(rctx.Param("sellerID"), 10, 64)
	if err != nil {
		res = fasthttp_tools.NewResponse_ErrBadRequest("sellerID is empty or wrong format")
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusBadRequest, *res)
		return nil
	}

	rating, err := c.service.Get(rctx, uint(sellerID))
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			res = fasthttp_tools.NewResponse_ErrNotFound("Seller id: " + strconv.Itoa(int(sellerID)) + " not found")
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		wblogger.Error(rctx, "service.Get error", err)
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	res = fasthttp_tools.NewResponse_Success(*rating)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}

func (c *blogController) MGet(rctx *routing.Context) error {
	const metricName = "blogController.MGet"
	var res *fasthttp_tools.Response

	sellerIDsStr := strings.Split(rctx.Param("sellerIDs"), ",")
	if len(sellerIDsStr) == 0 {
		res = fasthttp_tools.NewResponse_ErrBadRequest("sellerIDs is empty")
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusBadRequest, *res)
		return nil
	}

	sellerIDs := make([]uint, 0, len(sellerIDsStr))
	for _, s := range sellerIDsStr {
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			res = fasthttp_tools.NewResponse_ErrBadRequest("wrong format of sellerIDs")
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusBadRequest, *res)
			return nil
		}
		sellerIDs = append(sellerIDs, uint(i))
	}

	rating, err := c.service.MGet(rctx, &sellerIDs)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			res = fasthttp_tools.NewResponse_ErrNotFound("Sellers ids: " + fmt.Sprintf("%v", sellerIDs) + " not found")
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		wblogger.Error(rctx, "service.MGet error", err)
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	res = fasthttp_tools.NewResponse_Success(*rating)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}

func (c *blogController) Filter(rctx *routing.Context) error {
	var res *fasthttp_tools.Response
	args := rctx.URI().QueryArgs()

	filterParams := &blog.FilterParams{
		SellerOldIDs: string(args.Peek("sellerOldIds")),
		SortBy:       string(args.Peek("sortBy")),
		SortOrder:    string(args.Peek("sortOrder")),
		Limit:        string(args.Peek("limit")),
		Offset:       string(args.Peek("offset")),
	}

	if err := filterParams.Validate(); err != nil {
		res = fasthttp_tools.NewResponse_ErrBadRequest("filter params validation error. " + err.Error())
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
		return nil
	}

	filter, err := filterParams.Filter()
	if err != nil {
		res = fasthttp_tools.NewResponse_ErrBadRequest("filter params validation error. " + err.Error())
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
		return nil
	}

	count, ratings, err := c.service.Filter(rctx, filter)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			res = fasthttp_tools.NewResponse_ErrNotFound("")
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		wblogger.Error(rctx, "service.FilterInTsDB error", err)
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	result := NewResult_Ratings((*blog.Ratings)(ratings), *filter.Limit, *filter.Offset, count)
	res = fasthttp_tools.NewResponse_Success(*result)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}
*/
/*
func (c *blogController) FilterInTsDB(rctx *routing.Context) error {
	var res *fasthttp_tools.Response
	args := rctx.URI().QueryArgs()

	filterParams := &rating.FilterParams4TsDB{
		SellerOldIDs:       string(args.Peek("sellerOldIds")),
		SortBy:             string(args.Peek("sortBy")),
		SortOrder:          string(args.Peek("sortOrder")),
		Limit:              string(args.Peek("limit")),
		FromSellerOldID:    string(args.Peek("fromSellerOldId")),
		FromRating:         string(args.Peek("fromRating")),
		FromAvgBuyerRating: string(args.Peek("fromAvgBuyerRating")),
		FromRatioDefected:  string(args.Peek("fromRatioDefected")),
	}

	if err := filterParams.Validate(); err != nil {
		res = fasthttp_tools.NewResponse_ErrBadRequest("filter params validation error. " + err.Error())
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
		return nil
	}

	filter, err := filterParams.Filter()
	if err != nil {
		res = fasthttp_tools.NewResponse_ErrBadRequest("filter params validation error. " + err.Error())
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
		return nil
	}

	ratings, err := c.service.FilterInTsDB(rctx, filter)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			res = fasthttp_tools.NewResponse_ErrNotFound("")
			fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusNotFound, *res)
			return nil
		}
		wblogger.Error(rctx, "service.FilterInTsDB error", err)
		res = fasthttp_tools.NewResponse_ErrInternal()
		fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusInternalServerError, *res)
		return nil
	}

	res = fasthttp_tools.NewResponse_Success(*ratings)
	fasthttp_tools.FastHTTPWriteResult(rctx.RequestCtx, fasthttp.StatusOK, *res)
	return nil
}
*/

type result_blogs struct {
	Values     *[]blog.Blog               `json:"values"`
	Pagination *fasthttp_tools.Pagination `json:"pagination"`
}

func NewResult_Ratings(r *[]blog.Blog, limit uint, offset uint, count uint) *result_blogs {
	return &result_blogs{
		Values: r,
		Pagination: &fasthttp_tools.Pagination{
			Offset:  offset,
			Size:    limit,
			TotalNb: count,
		},
	}
}
