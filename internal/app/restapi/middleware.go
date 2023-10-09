package restapi

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/Kalinin-Andrey/blog/internal/pkg/fasthttp_tools"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/wildberries-tech/wblogger"
)

type httpMetrics interface {
	Inc(method, code, path, client string)
	WriteTiming(startTime time.Time, method, code, path, client string)
}

type Middlewares struct {
	metricsHandler fasthttp.RequestHandler
	metrics        httpMetrics
	metricsPath    []byte
	livePath       []byte
	readyPath      []byte
	authKey        *ecdsa.PublicKey
	authEnabled    bool
	ready          bool
}

func newMiddlewares() *Middlewares {
	return &Middlewares{}
}

func (m *Middlewares) Shutdown() {
	m.ready = false
}

func (m *Middlewares) WithMetrics(metrics httpMetrics, metricsPath []byte) *Middlewares {
	m.metricsPath = metricsPath
	m.metrics = metrics
	m.metricsHandler = fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	return m
}

func (m *Middlewares) WithProbes(livensess, readiness []byte) *Middlewares {
	m.livePath = livensess
	m.readyPath = readiness
	m.ready = true
	return m
}

func (m *Middlewares) WithAuth(key *ecdsa.PublicKey) *Middlewares {
	m.authKey = key
	m.authEnabled = true
	return m
}

func (m *Middlewares) RequestIdInterceptor(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValue(RequestIdKey, uuid.NewV4().String())
		handler(ctx)
	}
}

func (m *Middlewares) RecoverInterceptor(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				wblogger.Error(ctx, "PanicInterceptor", fmt.Errorf("%v", r))
				fasthttp_tools.InternalError(ctx, errors.New(fmt.Sprintf("%v", r)))
			}
		}()
		handler(ctx)
	}
}

func (m *Middlewares) LoggingInterceptor(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		wblogger.Debug(ctx, fmt.Sprintf("request: [%s] %s", ctx.Method(), string(ctx.URI().RequestURI())))
		next(ctx)
		wblogger.Debug(ctx, fmt.Sprintf("response: %d [%s]", ctx.Response.StatusCode(), ctx.Response.Body()))
	}
}

func (m *Middlewares) FastHTTPMetrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		now := time.Now().UTC()
		next(ctx)
		status := strconv.Itoa(ctx.Response.StatusCode())

		client := ""
		if s, ok := ctx.UserValue(AuthClientKey).(string); ok {
			client = s
		}

		path := "unknown"
		if s, ok := ctx.UserValue("metrics.label.path").(string); ok {
			path = s
		}

		m.metrics.Inc(string(ctx.Method()), status, path, client)
		m.metrics.WriteTiming(now, string(ctx.Method()), status, path, client)
	}
}

func (m *Middlewares) FastHTTPProbes(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if bytes.Equal(ctx.Request.URI().Path(), m.livePath) {
			ctx.SetStatusCode(http.StatusNoContent)
			return
		}

		if bytes.Equal(ctx.Request.URI().Path(), m.readyPath) {
			if m.ready {
				ctx.SetStatusCode(http.StatusNoContent)
				return
			}

			ctx.SetStatusCode(http.StatusInternalServerError)
			return
		}

		next(ctx)
	}
}

func (m *Middlewares) FastHTTPServeMetrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if bytes.Equal(ctx.Request.URI().Path(), m.metricsPath) {
			fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(ctx)
			return
		}
		next(ctx)
	}
}
