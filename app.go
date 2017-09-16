package gopress

import (
	"github.com/labstack/echo"
)

// App wrapper of echo.Echo and Container
type App struct {
	*echo.Echo

	Logger          *Logger
	Services        *Container
	MiddlewareGroup MiddlewareGroup
}

// MiddlewareGroup string key is the route prefix, value is []MiddlewareFunc
type MiddlewareGroup = map[string][]MiddlewareFunc

// AppContext is wrapper of echo.Context. It holds App instance of server.
type AppContext struct {
	echo.Context

	app *App
}

// App returns the App instance
func (c *AppContext) App() *App {
	return c.app
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

// GetGroup get a middlware group
func (app *App) GetRouteGroup(prefix string) *echo.Group {
	g := app.Group(prefix)
	if app.MiddlewareGroup == nil {
		return g
	}

	middlewares, ok := app.MiddlewareGroup[prefix]
	if !ok {
		return g
	}

	g.Use(middlewares...)

	return g
}

func (app *App) AddMiddlewareGroup(prefix string, m []MiddlewareFunc) {
	if app.MiddlewareGroup == nil {
		app.MiddlewareGroup = MiddlewareGroup{}
	}

	app.MiddlewareGroup[prefix] = m
}
