package middleWareFunc

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var RateLimiterConfig = middleware.RateLimiterConfig{
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: 60, Burst: 300, ExpiresIn: time.Minute},
	),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		return ctx.RealIP(), nil
	},
	ErrorHandler: func(context echo.Context, err error) error {
		return context.JSON(http.StatusTooManyRequests, nil)
	},
	DenyHandler: func(context echo.Context, identifier string, err error) error {
		return context.JSON(http.StatusForbidden, nil)
	},
}
