package gopress

import (
	"github.com/labstack/echo"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

const (
	requestIDContextKey = "request_id"
)

// App wrapper of echo.Echo and Container
type App struct {
	*echo.Echo

	Logger   *Logger
	Services *Container
}

// AppContext is wrapper of echo.Context. It holds App instance of server.
type AppContext struct {
	echo.Context

	app    *App
	logger *logrus.Entry
}

// App returns the App instance
func (c *AppContext) App() *App {
	return c.app
}

// Logger returns logger entry on current context
func (c *AppContext) RequestLogger() *logrus.Entry {
	return c.logger
}

// NewAppContextMiddleware returns a middleware which extends echo.Context
func NewAppContextMiddleware(app *App) MiddlewareFunc {
	return func(next HandlerFunc) echo.HandlerFunc {
		return func(c Context) error {

			// setup request id
			requestID := xid.New().String()
			c.Set(requestIDContextKey, requestID)
			logger := app.Logger.WithField("request_id", requestID)

			ac := &AppContext{c, app, logger}

			return next(ac)
		}
	}
}

// ContextApp try to get App instance from Context
func ContextApp(ctx Context) *App {
	ac, ok := ctx.(*AppContext)
	if !ok {
		return nil
	}
	return ac.App()
}

// ContextLogger returns logger entry for current request context
func ContextLogger(ctx Context) *logrus.Entry {
	if ctx, ok := ctx.(*AppContext); ok {
		return ctx.RequestLogger()
	}
	return defaultLogger.WithField("request_id", "")
}

// ContextRequestID returns ID for current request
func ContextRequestID(ctx Context) string {
	if id, ok := ctx.Get(requestIDContextKey).(string); ok {
		return id
	}
	return ""
}
