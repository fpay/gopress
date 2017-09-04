package gopress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// LoggingLevel is alias of logrus.Level
type LoggingLevel = logrus.Level

var (
	defaultLoggingFormatter = &logrus.JSONFormatter{}
	defaultLoggingLevel     = logrus.DebugLevel
	defaultLogger           = NewLogger()
)

// Logger wraps logrus.Logger
// TODO: implements echo.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger returns a Logger instance
func NewLogger() *Logger {
	l := &Logger{logrus.StandardLogger()}
	l.SetLevel(defaultLoggingLevel)
	l.SetOutput(os.Stdout)
	l.SetFormatter(defaultLoggingFormatter)
	return l
}

// SetOutput changes logger's output destination
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.Out = w
}

// SetFormatter changes logger's formatter
func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// NewLoggingMiddleware returns a middleware which logs every request
func NewLoggingMiddleware(name string, logger *Logger) MiddlewareFunc {

	// getLogger returns logrus.Logger. If user specify a logger when creating middleware, returns it.
	// If not, try to returns App's logger. If app is not found on the context, returns the standard logger.
	getLogger := func(c Context) *Logger {
		if logger != nil {
			return logger
		}

		if app := AppFromContext(c); app != nil {
			return app.Logger
		}

		return defaultLogger
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			l := getLogger(c)
			start := time.Now()

			req := c.Request()
			entry := l.WithFields(logrus.Fields{
				"host":     req.Host,
				"remote":   req.RemoteAddr,
				"method":   req.Method,
				"uri":      req.RequestURI,
				"referer":  req.Referer(),
				"bytes_in": req.ContentLength,
				"scope":    name,
			})

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)

			resp := c.Response()
			entry.WithFields(logrus.Fields{
				"status":    resp.Status,
				"bytes_out": resp.Size,
				"latency":   fmt.Sprintf("%.3f", latency.Seconds()*1000),
			}).Info("request completes.")

			return nil
		}
	}
}
