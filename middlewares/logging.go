package middlewares

import (
	"fmt"
	"time"

	"github.com/fpay/gopress/log"
	"github.com/fpay/gopress/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type LoggingMiddlewareConfig struct {
	Name    string
	Skipper middleware.Skipper
}

// NewLoggingMiddleware returns a middleware which logs every request
func NewLoggingMiddleware(opts LoggingMiddlewareConfig) echo.MiddlewareFunc {

	name := opts.Name
	skipper := opts.Skipper

	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			l := log.Extract(c)
			start := time.Now()

			req := c.Request()

			httpFields := logrus.Fields{
				"host":     req.Host,
				"remote":   utils.RequestRemoteAddr(req),
				"method":   req.Method,
				"uri":      req.RequestURI,
				"referer":  req.Referer(),
				"bytes_in": req.ContentLength,
			}
			entry := l.WithFields(logrus.Fields{
				"scope": name,
				"http":  httpFields,
			})

			if err := next(c); err != nil {
				c.Error(err)
				entry = entry.WithError(err)
			}

			latency := time.Since(start)

			resp := c.Response()

			httpFields["status"] = resp.Status
			httpFields["bytes_out"] = resp.Size
			httpFields["latency"] = fmt.Sprintf("%.3f", latency.Seconds()*1000)
			entry.WithFields(logrus.Fields{
				"http": httpFields,
			}).Info("request completes.")

			return nil
		}
	}
}
