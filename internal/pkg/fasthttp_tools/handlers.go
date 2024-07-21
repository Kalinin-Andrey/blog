package fasthttp_tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	uuid "github.com/satori/go.uuid"

	"github.com/valyala/fasthttp"
)

const (
	UserCtxField  = "logger.userId"
	TxIdCtxField  = "logger.txId"
	SumCtxField   = "logger.sum"
	StateCtxField = "logger.state"

	AuthClientKey = "http.client"
	RequestIdKey  = "http.requestId"
)

var stopList = map[string]struct{}{}

func parseUserID(ctx *fasthttp.RequestCtx) (string, error) {
	v, ok := ctx.UserValue("user_id").(string)
	if !ok {
		return "", errors.New("invalid userId")
	}

	if len(v) > 20 {
		return "", errors.New("invalid userId")
	}

	ctx.SetUserValue(UserCtxField, v)

	if _, ok := stopList[v]; ok {
		return "", errors.New("stoplist")
	}

	return v, nil
}

func parsePaymentID(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	v, ok := ctx.UserValue("payment_id").(string)
	if !ok {
		return uuid.Nil, errors.New("invalid payment_id")
	}
	return uuid.FromString(v)
}

func parseNonNegativeSum(ctx *fasthttp.RequestCtx) (int64, error) {
	sumV := string(ctx.QueryArgs().Peek("sum"))
	sum, err := strconv.ParseInt(sumV, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid sum value: %w", err)
	}

	if sum < 0 {
		return 0, errors.New("negative sum")
	}

	ctx.SetUserValue(SumCtxField, sumV)

	return sum, nil
}

func parseTxID(ctx *fasthttp.RequestCtx) (string, error) {
	txId := string(ctx.QueryArgs().Peek("tx_id"))
	if txId == "" {
		return "", fmt.Errorf("invalid txId value: %s", txId)
	}

	ctx.SetUserValue(TxIdCtxField, txId)

	return txId, nil
}

func BadRequest(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	resp := errorResp{
		Error:     err.Error(),
		RequestId: "",
	}
	if id, ok := ctx.UserValue(RequestIdKey).(string); ok {
		resp.RequestId = id
	}

	res, _ := json.Marshal(resp)

	ctx.Write(res)
}

func InternalError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	resp := errorResp{
		Error:     err.Error(),
		RequestId: "",
	}
	if id, ok := ctx.UserValue(RequestIdKey).(string); ok {
		resp.RequestId = id
	}

	res, _ := json.Marshal(resp)
	ctx.Write(res)
}

type errorResp struct {
	Error     string
	RequestId string
}

func Success(ctx *fasthttp.RequestCtx, body []byte) error {
	if body == nil {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil
	}
	if _, err := ctx.Write(body); err != nil {
		return fmt.Errorf("Write(body) error: %w", err)
	}
	return nil
}

func parseUserId(ctx *fasthttp.RequestCtx) (string, error) {
	v, ok := ctx.UserValue("user_id").(string)
	if !ok {
		return "", errors.New("invalid userId")
	}

	if len(v) > 20 {
		return "", errors.New("invalid userId")
	}

	return v, nil
}
