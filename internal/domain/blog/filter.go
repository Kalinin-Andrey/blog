package blog

import (
	"fmt"
	"github.com/Kalinin-Andrey/blog/internal/pkg/apperror"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"regexp"
	"strconv"
	"strings"
)

const (
	SortOrderAsc                = "asc"
	SortOrderDesc               = "desc"
	JsonProperty_SellerOldID    = "sellerOldId"
	JsonProperty_Rating         = "rating"
	JsonProperty_AvgBuyerRating = "avgBuyerRating"
	JsonProperty_RatioDefected  = "ratioDefected"
	Default_SortBy              = JsonProperty_SellerOldID
	Default_SortOrder           = SortOrderAsc
	Default_Limit               = 100
	Default_Offset              = 0
)

var possiblePropertiesForSorting = []interface{}{
	JsonProperty_SellerOldID,
	JsonProperty_Rating,
	JsonProperty_AvgBuyerRating,
	JsonProperty_RatioDefected,
}

var possibleSortOrders = []interface{}{
	SortOrderAsc,
	SortOrderDesc,
}

type Filter struct {
	SellerOldIDs    *[]uint
	SortBy          *string
	SortOrder       *string
	IsSortOrderDesc bool
	Limit           *uint
	Offset          *uint
}

func NewFilter(sellerOldIDs *[]uint, sortBy *string, sortOrder *string, limit *uint, offset *uint) (*Filter, error) {
	f := &Filter{
		SellerOldIDs: sellerOldIDs,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Limit:        limit,
		Offset:       offset,
	}

	if f.SortBy == nil || *f.SortBy == "" {
		s := Default_SortBy
		f.SortBy = &s
	}
	if f.SortOrder == nil || *f.SortOrder == "" {
		s := Default_SortOrder
		f.SortOrder = &s
	}
	if *f.SortOrder == SortOrderDesc {
		f.IsSortOrderDesc = true
	}
	if f.Limit == nil || *f.Limit == 0 {
		u := uint(Default_Limit)
		f.Limit = &u
	}
	if f.Offset == nil {
		u := uint(Default_Offset)
		f.Offset = &u
	}
	return f, f.Validate()
}

func (e *Filter) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.SellerOldIDs, validation.NilOrNotEmpty),
		validation.Field(&e.SortBy, validation.Required, validation.In(possiblePropertiesForSorting...)),
		validation.Field(&e.SortOrder, validation.Required, validation.In(possibleSortOrders...)),
	)
}

type FilterParams struct {
	SellerOldIDs string `json:"sellerOldIDs"`
	SortBy       string `json:"sortBy"`
	SortOrder    string `json:"sortOrder"`
	Limit        string `json:"limit"`
	Offset       string `json:"offset"`
}

func (e *FilterParams) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.SellerOldIDs, validation.Match(regexp.MustCompile("^[0-9,]+$"))),
		validation.Field(&e.SortBy, validation.In(possiblePropertiesForSorting...)),
		validation.Field(&e.SortOrder, validation.In(possibleSortOrders...)),
		validation.Field(&e.Limit, is.Int),
		validation.Field(&e.Offset, is.Int),
	)
}

func (e *FilterParams) Filter() (*Filter, error) {
	f := Filter{}

	if e.SellerOldIDs != "" {
		sellerIDsStr := strings.Split(e.SellerOldIDs, ",")
		sellerOldIDs := make([]uint, 0, len(sellerIDsStr))
		for _, s := range sellerIDsStr {
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("[%w] sellerOldIDs parse error: %s", apperror.ErrBadRequest, err.Error())
			}
			sellerOldIDs = append(sellerOldIDs, uint(i))
		}
		f.SellerOldIDs = &sellerOldIDs
	}

	if e.SortBy != "" {
		f.SortBy = &e.SortBy
	}

	if e.SortOrder != "" {
		f.SortOrder = &e.SortOrder
	}

	if e.Limit != "" {
		u, err := strconv.ParseUint(e.Limit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] limit parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		limit := uint(u)
		f.Limit = &limit
	}

	if e.Offset != "" {
		u, err := strconv.ParseUint(e.Offset, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] offset parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		offset := uint(u)
		f.Offset = &offset
	}

	return NewFilter(f.SellerOldIDs, f.SortBy, f.SortOrder, f.Limit, f.Offset)
}

type Filter4TsDB struct {
	SellerOldIDs       *[]uint
	SortBy             *string
	SortOrder          *string
	IsSortOrderDesc    bool
	Limit              *uint
	FromSellerOldID    *uint
	FromRating         *float64
	FromAvgBuyerRating *float64
	FromRatioDefected  *float64
}

func NewFilter4TsDB(SellerOldIDs *[]uint, SortBy *string, SortOrder *string, Limit *uint, FromSellerOldID *uint, FromRating *float64, FromAvgBuyerRating *float64, FromRatioDefected *float64) (*Filter4TsDB, error) {
	f := &Filter4TsDB{
		SellerOldIDs:       SellerOldIDs,
		SortBy:             SortBy,
		SortOrder:          SortOrder,
		Limit:              Limit,
		FromSellerOldID:    FromSellerOldID,
		FromRating:         FromRating,
		FromAvgBuyerRating: FromAvgBuyerRating,
		FromRatioDefected:  FromRatioDefected,
	}

	if f.SortBy == nil || *f.SortBy == "" {
		s := Default_SortBy
		f.SortBy = &s
	}
	if f.SortOrder == nil || *f.SortOrder == "" {
		s := Default_SortOrder
		f.SortOrder = &s
	}
	if *f.SortOrder == SortOrderDesc {
		f.IsSortOrderDesc = true
	}
	if f.Limit == nil || *f.Limit == 0 {
		u := uint(Default_Limit)
		f.Limit = &u
	}
	return f, f.Validate()
}

func (e *Filter4TsDB) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.SellerOldIDs, validation.NilOrNotEmpty),
		validation.Field(&e.SortBy, validation.Required, validation.In(possiblePropertiesForSorting...)),
		validation.Field(&e.SortOrder, validation.Required, validation.In(possibleSortOrders...)),
		//validation.Field(&e.Limit, validation.Required),
		validation.Field(&e.FromSellerOldID, validation.When(e.FromRating != nil, validation.Required), validation.When(e.FromAvgBuyerRating != nil, validation.Required), validation.When(e.FromRatioDefected != nil, validation.Required)),
	)
}

type FilterParams4TsDB struct {
	SellerOldIDs       string `json:"sellerOldIDs"`
	SortBy             string `json:"sortBy"`
	SortOrder          string `json:"sortOrder"`
	Limit              string `json:"limit"`
	FromSellerOldID    string `json:"fromSellerOldID"`
	FromRating         string `json:"fromRating"`
	FromAvgBuyerRating string `json:"fromAvgBuyerRating"`
	FromRatioDefected  string `json:"fromRatioDefected"`
}

func (e *FilterParams4TsDB) Validate() error {
	return validation.ValidateStruct(e,
		validation.Field(&e.SellerOldIDs, validation.Match(regexp.MustCompile("^[0-9,]+$"))),
		validation.Field(&e.SortBy, validation.In(possiblePropertiesForSorting...)),
		validation.Field(&e.SortOrder, validation.In(possibleSortOrders...)),
		validation.Field(&e.Limit, is.Int),
		validation.Field(&e.FromSellerOldID, is.Int, validation.When(e.FromRating != "", validation.Required), validation.When(e.FromAvgBuyerRating != "", validation.Required), validation.When(e.FromRatioDefected != "", validation.Required)),
		validation.Field(&e.FromRating, is.Float),
		validation.Field(&e.FromAvgBuyerRating, is.Float),
		validation.Field(&e.FromRatioDefected, is.Float),
	)
}

func (e *FilterParams4TsDB) Filter() (*Filter4TsDB, error) {
	f := Filter4TsDB{}

	if e.SellerOldIDs != "" {
		sellerIDsStr := strings.Split(e.SellerOldIDs, ",")
		sellerOldIDs := make([]uint, 0, len(sellerIDsStr))
		for _, s := range sellerIDsStr {
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("[%w] sellerOldIDs parse error: %s", apperror.ErrBadRequest, err.Error())
			}
			sellerOldIDs = append(sellerOldIDs, uint(i))
		}
		f.SellerOldIDs = &sellerOldIDs
	}

	if e.SortBy != "" {
		f.SortBy = &e.SortBy
	}

	if e.SortOrder != "" {
		f.SortOrder = &e.SortOrder
	}

	if e.Limit != "" {
		u, err := strconv.ParseUint(e.Limit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] limit parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		limit := uint(u)
		f.Limit = &limit
	}

	if e.FromSellerOldID != "" {
		u, err := strconv.ParseUint(e.FromSellerOldID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] fromSellerOldID parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		fromSellerOldID := uint(u)
		f.FromSellerOldID = &fromSellerOldID
	}

	if e.FromRating != "" {
		fl, err := strconv.ParseFloat(e.FromRating, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] fromRating parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		f.FromRating = &fl
	}

	if e.FromAvgBuyerRating != "" {
		fl, err := strconv.ParseFloat(e.FromAvgBuyerRating, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] fromAvgBuyerRating parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		f.FromAvgBuyerRating = &fl
	}

	if e.FromRatioDefected != "" {
		fl, err := strconv.ParseFloat(e.FromRatioDefected, 64)
		if err != nil {
			return nil, fmt.Errorf("[%w] fromRatioDefected parse error: %s", apperror.ErrBadRequest, err.Error())
		}
		f.FromRatioDefected = &fl
	}

	return NewFilter4TsDB(f.SellerOldIDs, f.SortBy, f.SortOrder, f.Limit, f.FromSellerOldID, f.FromRating, f.FromAvgBuyerRating, f.FromRatioDefected)
}
