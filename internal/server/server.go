package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapiecho "github.com/oapi-codegen/echo-middleware"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
	"github.com/ziflex/rm-rf-production/internal/api"
)

type (
	Server struct {
		engine *echo.Echo
	}

	Options struct {
		Logger zerolog.Logger
		Spec   []byte
		UI     fs.FS
	}
)

func NewServer(handler api.StrictServerInterface, opts Options) (*Server, error) {
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromData(opts.Spec)

	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	echoLogger := lecho.From(opts.Logger)
	svr := &Server{}
	svr.engine = echo.New()
	svr.engine.Logger = echoLogger
	svr.engine.HideBanner = true
	svr.engine.HTTPErrorHandler = errorHandler

	svr.engine.Use(middleware.BodyLimit("1M"))
	svr.engine.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		TargetHeader: echo.HeaderXCorrelationID,
	}))
	svr.engine.Use(lecho.Middleware(lecho.Config{
		Logger:          echoLogger,
		RequestIDKey:    "request_id",
		RequestIDHeader: echo.HeaderXCorrelationID,
		HandleError:     false,
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/health" || c.Path() == "/openapi.yaml" || strings.HasPrefix(c.Path(), "/docs")
		},
	}))
	spec.Servers = nil
	svr.engine.Use(oapiecho.OapiRequestValidatorWithOptions(spec, &oapiecho.Options{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path

			return path == "/health" || path == "/openapi.yaml" || strings.HasPrefix(path, "/docs")
		},
	}))
	svr.engine.Use(middleware.Recover())
	svr.engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	svr.engine.GET("/openapi.yaml", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/x-yaml", opts.Spec)
	})

	if opts.UI != nil {
		svr.engine.StaticFS("/docs", opts.UI)
	}

	api.RegisterHandlers(svr.engine, api.NewStrictHandler(handler, nil))

	return svr, nil
}

func (svr *Server) Run(port int) error {
	return svr.engine.Start(fmt.Sprintf("0.0.0.0:%d", port))
}

func (svr *Server) Shutdown() error {
	return svr.engine.Shutdown(nil)
}
