package gopress

import (
	"github.com/labstack/echo"
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

	app *App
}

// App returns the App instance
func (c *AppContext) App() *App {
	return c.app
}

// NewAppContextMiddleware returns a middleware which extends echo.Context
func NewAppContextMiddleware(app *App) MiddlewareFunc {
	return func(next HandlerFunc) echo.HandlerFunc {
		return func(c Context) error {
			ac := &AppContext{c, app}
			return next(ac)
		}
	}
}

// AppFromContext try to get App instance from Context
func AppFromContext(ctx Context) *App {
	ac, ok := ctx.(*AppContext)
	if !ok {
		return nil
	}
	return ac.App()
}
