package gopress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var (
	defaultLoggingOutput    = os.Stdout
	defaultLoggingFormatter = &logrus.JSONFormatter{}
	defaultLoggingLevel     = log.DEBUG
	defaultLogger           = NewLogger()
	loggingLevelMapping     = map[log.Lvl]logrus.Level{
		log.DEBUG: logrus.DebugLevel,
		log.INFO:  logrus.InfoLevel,
		log.WARN:  logrus.WarnLevel,
		log.ERROR: logrus.ErrorLevel,
		log.OFF:   logrus.PanicLevel,
	}
)

// Logger wraps logrus.Logger
type Logger struct {
	*logrus.Logger

	level log.Lvl
}

// NewLogger returns a Logger instance
func NewLogger() *Logger {
	l := &Logger{Logger: &logrus.Logger{}}
	l.SetLevel(defaultLoggingLevel)
	l.SetOutput(os.Stdout)
	l.SetFormatter(defaultLoggingFormatter)
	return l
}

// Output returns Logger's output destination.
func (l *Logger) Output() io.Writer {
	return l.Logger.Out
}

// SetOutput changes logger's output destination
func (l *Logger) SetOutput(w io.Writer) {
	l.Logger.Out = w
}

// SetFormatter changes logger's formatter
func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// Prefix is used to implement echo.Logger.
// This function always returns empty string because prefix is not needed in logrus.
func (l *Logger) Prefix() string {
	return ""
}

// SetPrefix is used to implement echo.Logger. Do nothing here.
func (l *Logger) SetPrefix(p string) {}

// Level returns current logging level.
func (l *Logger) Level() log.Lvl {
	return l.level
}

// SetLevel changes logging level.
// If you want to change underlying logrus logger's level, call l.Logger.SetLevel function.
func (l *Logger) SetLevel(v log.Lvl) {
	l.level = v
	l.Logger.SetLevel(loggingLevelMapping[v])
}

// Printj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Print.
func (l *Logger) Printj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Print()
}

// Debugj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Debug.
func (l *Logger) Debugj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Debug()
}

// Infoj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Info.
func (l *Logger) Infoj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Info()
}

// Warnj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Warn.
func (l *Logger) Warnj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Warn()
}

// Errorj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Error.
func (l *Logger) Errorj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Error()
}

// Fatalj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Fatal.
func (l *Logger) Fatalj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Fatal()
}

// Panicj is used to implement echo.Logger. It creates an logrus.Entry with fields j, then call Panic.
func (l *Logger) Panicj(j log.JSON) {
	l.WithFields(logrus.Fields(j)).Panic()
}

// NewLoggingMiddleware returns a middleware which logs every request
func NewLoggingMiddleware(name string, logger *Logger) MiddlewareFunc {

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
			l := getLogger(c)
			start := time.Now()

			req := c.Request()
			entry := l.WithFields(logrus.Fields{
				"http.host":     req.Host,
				"http.remote":   RequestRemoteAddr(req),
				"http.method":   req.Method,
				"http.uri":      req.RequestURI,
				"http.referer":  req.Referer(),
				"http.bytes_in": req.ContentLength,
				"scope":         name,
			})

			if err := next(c); err != nil {
				c.Error(err)
				entry = entry.WithError(err)
			}

			latency := time.Since(start)

			resp := c.Response()
			entry.WithFields(logrus.Fields{
				"http.status":    resp.Status,
				"http.bytes_out": resp.Size,
				"http.latency":   fmt.Sprintf("%.3f", latency.Seconds()*1000),
			}).Info("request completes.")

			return nil
		}
	}
}
