package gopress

import (
	"github.com/fpay/gopress/log"
	"github.com/labstack/echo/v4"
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
}

// AppContext is wrapper of echo.Context. It holds App instance of server.
type AppContext struct {
	echo.Context

	app    *App
}

func NewAppContext(c echo.Context, app *App, logger *logrus.Entry) *AppContext {
	return &AppContext{c, app}
}

// App returns the App instance
func (c *AppContext) App() *App {
	return c.app
}

// NewAppContextMiddleware returns a middleware which extends echo.Context
func NewAppContextMiddleware(app *App) MiddlewareFunc {
	return func(next HandlerFunc) echo.HandlerFunc {
		return func(c Context) error {

			// setup request id
			requestID := xid.New().String()
			c.Set(requestIDContextKey, requestID)

			ac := &AppContext{c, app}

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
var ContextLogger = log.Extract

var ContextRequestID = ExtractRequestID

// ExtractRequestID returns ID for current request
func ExtractRequestID(ctx Context) string {
	if id, ok := ctx.Get(requestIDContextKey).(string); ok {
		return id
	}
	return ""
}
