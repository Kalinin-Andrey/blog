package fasthttp_tools

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func FastHTTPWriteResult(ctx *fasthttp.RequestCtx, status int, data interface{}) error {

	if data != nil {
		var body []byte
		const metricName = "FastHTTPWriteResult "

		body, err := json.Marshal(data)
		if err != nil {
			status = fasthttp.StatusInternalServerError
			data = NewResponse_ErrInternal()
			body, _ = json.Marshal(data)
			return fmt.Errorf(metricName+"json.Marshal error: %w", err)
		}

		if _, err := ctx.Write(body); err != nil {
			return fmt.Errorf("Write(body) error: %w", err)
		}
	}
	ctx.SetStatusCode(status)

	return nil
}
