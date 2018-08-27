package gopress

import (
	"fmt"
	"time"

	"github.com/fpay/gopress/log"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger
type Logger = log.Logger

var defaultLogger = log.DefaultLogger

type LoggingMiddlewareConfig struct {
	Name    string
	Logger  *Logger
	Skipper middleware.Skipper
}

// NewLoggingMiddleware returns a middleware which logs every request
func NewLoggingMiddleware(opts LoggingMiddlewareConfig) MiddlewareFunc {

	name := opts.Name
	logger := opts.Logger
	skipper := opts.Skipper

	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	// getLogger returns Logger. If user specify a logger when creating middleware, returns it.
	// If not, try to returns App's logger. If app is not found on the context, returns the default logger.
	getLogger := func(c Context) *logrus.Entry {
		if logger != nil {
			return logger.WithFields(nil)
		}

		return ContextLogger(c)
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			if skipper(c) {
				return next(c)
			}

			l := getLogger(c)
			start := time.Now()

			req := c.Request()

			httpFields := logrus.Fields{
				"host":     req.Host,
				"remote":   RequestRemoteAddr(req),
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
