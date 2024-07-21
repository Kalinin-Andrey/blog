package fasthttp_tools

import "blog/internal/pkg/apperror"

type Response struct {
	Data             interface{} `json:"data"`
	Error            bool        `json:"error"`
	ErrorText        error       `json:"errorText"`
	AdditionalErrors error       `json:"additionalErrors"`
	Pagination       *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Offset  uint `json:"offset"`
	Size    uint `json:"size"`
	TotalNb uint `json:"totalNb"`
}

func NewResponse_ErrUnauthorized() *Response {
	return &Response{
		Error:     true,
		ErrorText: apperror.NewError(""),
	}
}

func NewResponse_ErrForbidden() *Response {
	return &Response{
		Error:     true,
		ErrorText: apperror.NewError(""),
	}
}

func NewResponse_ErrBadRequest(errMessage string) *Response {
	return &Response{
		Error:     true,
		ErrorText: apperror.NewError(errMessage),
	}
}

func NewResponse_ErrNotFound(errMessage string) *Response {
	if errMessage == "" {
		errMessage = "not found"
	}
	return &Response{
		Error:     true,
		ErrorText: apperror.NewError(errMessage),
	}
}

func NewResponse_ErrInternal() *Response {
	return &Response{
		Error:     true,
		ErrorText: apperror.NewError("internal error"),
	}
}

func NewResponse_Success(data interface{}) *Response {
	return &Response{
		ErrorText: apperror.NewError(""),
		Data:      data,
	}
}

//
//func NewResponse_SuccessWithPagination(data interface{}, limit uint, offset uint, count uint) *Response {
//	return &Response{
//		ErrorText: apperror.NewError(""),
//		Data:      data,
//		Pagination: &Pagination{
//			Offset:  offset,
//			Size:    limit,
//			TotalNb: count,
//		},
//	}
//}
