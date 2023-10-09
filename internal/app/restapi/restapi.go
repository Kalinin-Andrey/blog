package restapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"github.com/wildberries-tech/wblogger"

	"github.com/Kalinin-Andrey/blog/internal/pkg/config"
	"github.com/Kalinin-Andrey/blog/internal/pkg/fasthttp_tools"

	"github.com/Kalinin-Andrey/blog/internal/app"
	"github.com/minipkg/prometheus-utils"
)

const (
	AuthClientKey = "http.client"
	RequestIdKey  = "X-Request-Id"
)

type HttpServerMetric interface {
	Inc(method, code, path, client string)
	WriteTiming(startTime time.Time, method, code, path, client string)
}

type RestAPI struct {
	*app.App
	config              *config.API
	serverRestAPI       *fasthttp.Server
	serverRestAPIMetric HttpServerMetric
	serverMetrics       *fasthttp.Server
	serverProbes        *fasthttp.Server
}

func New(app *app.App, appConfig *config.AppConfig, cfg *config.API) *RestAPI {
	restAPI := &RestAPI{
		config: cfg,
		App:    app,
		serverRestAPI: &fasthttp.Server{
			Name:            appConfig.Name,
			ReadTimeout:     cfg.Rest.ReadTimeout,
			WriteTimeout:    cfg.Rest.WriteTimeout,
			IdleTimeout:     cfg.Rest.IdleTimeout,
			CloseOnShutdown: true,
		},
		serverRestAPIMetric: prometheus_utils.NewHttpServerMetrics(appConfig.NameSpace, appConfig.Name, appConfig.Service),
		serverMetrics: &fasthttp.Server{
			Name:            appConfig.Name,
			ReadTimeout:     cfg.Metrics.ReadTimeout,
			WriteTimeout:    cfg.Metrics.WriteTimeout,
			IdleTimeout:     cfg.Metrics.IdleTimeout,
			CloseOnShutdown: true,
		},
		serverProbes: &fasthttp.Server{
			Name:            appConfig.Name,
			ReadTimeout:     cfg.Probes.ReadTimeout,
			WriteTimeout:    cfg.Probes.WriteTimeout,
			IdleTimeout:     cfg.Probes.IdleTimeout,
			CloseOnShutdown: true,
		},
	}

	wblogger.CtxField(AuthClientKey)
	wblogger.CtxField(RequestIdKey)
	wblogger.CtxField(fasthttp_tools.SumCtxField)
	wblogger.CtxField(fasthttp_tools.TxIdCtxField)
	wblogger.CtxField(fasthttp_tools.UserCtxField)

	restAPI.buildHandler()

	return restAPI
}

func (a *RestAPI) buildHandler() {
	rp := routing.New()
	rp.Get("/live", LiveHandler)
	rp.Get("/ready", LiveHandler)
	a.serverProbes.Handler = rp.HandleRequest

	rm := routing.New()
	rm.Get("/metrics", prometheus_utils.GetFasthttpRoutingHandler())
	a.serverMetrics.Handler = rm.HandleRequest

	r := routing.New()
	r.Use(RecoverInterceptorMiddleware, SetResponseHeaderMiddleware("Content-Type", "application/json; charset=utf-8"), RequestIdInterceptorMiddleware, a.httpServerMetricMiddleware)
	//api := r.Group("/api/v1")

	//ratingController := controller.NewBlogController(r, a.Domain.Blog)
	//api.Get("/rating/<sellerID>", ratingController.Get)
	a.serverRestAPI.Handler = r.HandleRequest

}

func SetResponseHeaderMiddleware(key string, value string) func(rctx *routing.Context) error {
	return func(rctx *routing.Context) error {
		rctx.Response.Header.Set(key, value)
		return rctx.Next()
	}
}

func RecoverInterceptorMiddleware(rctx *routing.Context) error {
	defer func() {
		if r := recover(); r != nil {
			wblogger.Error(rctx, "PanicInterceptor", fmt.Errorf("%v", r))
			fasthttp_tools.InternalError(rctx.RequestCtx, fmt.Errorf("%v", r))
		}
	}()

	return rctx.Next()
}

func RequestIdInterceptorMiddleware(rctx *routing.Context) error {
	if requestId := rctx.Get(RequestIdKey); requestId != nil {
		return nil
	}
	if requestIdB := rctx.RequestCtx.Request.Header.Peek(RequestIdKey); requestIdB != nil && len(requestIdB) > 0 {
		rctx.Set(RequestIdKey, string(requestIdB))
		return nil
	}
	rctx.Set(RequestIdKey, uuid.NewV4().String())

	// пока больше клиентов нет, потом - посмотрим
	rctx.Set(AuthClientKey, "admin")

	return rctx.Next()
}

func (a *RestAPI) httpServerMetricMiddleware(rctx *routing.Context) error {
	now := time.Now()

	err := rctx.Next()

	client := ""
	if s, ok := rctx.Get(AuthClientKey).(string); ok {
		client = s
	}

	method := string(rctx.Method())
	path := string(rctx.Path())
	status := strconv.Itoa(rctx.Response.StatusCode())

	a.serverRestAPIMetric.Inc(method, status, path, client)
	a.serverRestAPIMetric.WriteTiming(now, method, status, path, client)

	return err
}

func LiveHandler(rctx *routing.Context) error {
	rctx.SetStatusCode(http.StatusNoContent)
	return nil
}

func (a *RestAPI) Run(ctx context.Context) error {
	if err := a.App.Run(); err != nil {
		return err
	}
	go func() {
		wblogger.Info(context.Background(), "metrics listen on "+a.config.Metrics.Addr)
		if err := a.serverMetrics.ListenAndServe(a.config.Metrics.Addr); err != nil {
			wblogger.Error(ctx, "serverMetrics.ListenAndServe error", err)
			wblogger.Flush()
		}
	}()
	go func() {
		wblogger.Info(context.Background(), "probes listen on "+a.config.Probes.Addr)
		if err := a.serverProbes.ListenAndServe(a.config.Probes.Addr); err != nil {
			wblogger.Error(ctx, "serverProbes.ListenAndServe error", err)
			wblogger.Flush()
		}
	}()
	wblogger.Info(context.Background(), "restapi listen on "+a.config.Rest.Addr)
	return a.serverRestAPI.ListenAndServe(a.config.Rest.Addr)
}

func (a *RestAPI) Stop() error {
	var err error
	var errs []error

	wblogger.Info(context.Background(), "Api-Shutdown")
	time.Sleep(time.Second * 10)

	if err = a.serverRestAPI.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err = a.serverMetrics.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err = a.serverProbes.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err = a.App.Stop(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
