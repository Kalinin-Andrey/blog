package fasthttp_tools

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
	"github.com/wildberries-tech/wblogger"
)

func FastHTTPWriteResult(ctx *fasthttp.RequestCtx, status int, data interface{}) {

	if data != nil {
		var body []byte
		const metricName = "FastHTTPWriteResult "

		body, err := json.Marshal(data)
		if err != nil {
			status = fasthttp.StatusInternalServerError
			data = NewResponse_ErrInternal()
			body, _ = json.Marshal(data)
			wblogger.Error(ctx, metricName+"json.Marshal error: ", err)
		}

		if _, err := ctx.Write(body); err != nil {
			wblogger.Error(ctx, metricName+"WriteResp-Err", err)
		}
	}
	ctx.SetStatusCode(status)

	return
}
