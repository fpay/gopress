package gopress

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// App wrapper of echo.Echo and Container
type App struct {
	*echo.Echo

	Logger   *logrus.Logger
	Services *Container
}

// App returns App instance of server
func (s *Server) App() *App {
	return s.app
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

// Services returns Container of current App
func (c *AppContext) Services() *Container {
	return c.app.Services
}

// appContextMiddleware returns a middleware which extends echo.Context
func appContextMiddleware(app *App) MiddlewareFunc {
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
