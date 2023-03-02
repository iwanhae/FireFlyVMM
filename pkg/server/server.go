package server

import (
	"net/http"
	"time"

	"github.com/iwanhae/monolithcloud/pkg/vmm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

type ServerOpts struct {
	VMM *vmm.VirtualMachineManager
}

func NewServer(opts ServerOpts) http.Handler {
	e := echo.New()

	// Logger
	e.Use(
		echo.WrapMiddleware(hlog.NewHandler(log.Logger)),
		echo.WrapMiddleware(hlog.AccessHandler(
			func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Stringer("url", r.URL).
					Int("status", status).
					Dur("duration", duration).
					Send()
			})),
		echo.WrapMiddleware(hlog.RemoteAddrHandler("ip")),
		echo.WrapMiddleware(hlog.UserAgentHandler("user_agent")),
		echo.WrapMiddleware(hlog.RequestIDHandler("request_id", "Request-Id")),
	)

	// etc
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				c.JSON(500, Error{
					Code:    "-1",
					Message: err.Error(),
				})
			}
			return nil
		}
	})
	e.Use(middleware.Recover())

	// API
	RegisterHandlers(e, &Server{
		vmm: opts.VMM,
	})

	return e
}
