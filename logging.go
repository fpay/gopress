package gopress

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// LoggingLevel is alias of logrus.Level
type LoggingLevel = log.Level

var (
	defaultLoggingFormatter = &log.JSONFormatter{}
	defaultLoggingLevel     = log.DebugLevel
)

// init set defaults for logging
func init() {
	log.SetFormatter(defaultLoggingFormatter)
	log.SetLevel(defaultLoggingLevel)
	log.SetOutput(os.Stdout)
}

// SetLoggingLevel changes logging level
func (s *Server) SetLoggingLevel(level LoggingLevel) {
	log.SetLevel(level)
}

// SetLoggingOutput changes logging output destination
func (s *Server) SetLoggingOutput(w io.Writer) {
	log.SetOutput(w)
}

// NewLoggingMiddleware returns logging middleware handler function
func NewLoggingMiddleware(name string) MiddlewareFunc {
	l := log.StandardLogger()
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			start := time.Now()

			req := c.Request()
			entry := l.WithFields(log.Fields{
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
			entry.WithFields(log.Fields{
				"status":    resp.Status,
				"bytes_out": resp.Size,
				"latency":   fmt.Sprintf("%.3f", latency.Seconds()*1000),
			}).Info("request completes.")

			return nil
		}
	}
}
