package gopress

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
)

func TestContextApp(t *testing.T) {
	app := &App{}
	c := &AppContext{app: app}
	actual := c.App()
	if actual != app {
		t.Errorf("expect app is %#v, actual is %#v", app, actual)
	}
}

func TestAppContextMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	app := &App{}
	var actual Context
	h := appContextMiddleware(app)(func(c Context) error {
		actual = c
		return c.String(http.StatusOK, "test context")
	})
	h(c)

	if c, ok := actual.(*AppContext); !ok {
		t.Errorf("expect context is AppContext, actual is %#v", c)
	}
}

func TestAppAddMiddlewareGroup(t *testing.T) {
	app := &App{}
	app.Echo = echo.New()

	m := make([]MiddlewareFunc, 0)
	m = append(m, NewLoggingMiddleware("global", NewLogger()))
	app.AddMiddlewareGroup("/hello", m)
	m2, ok := app.MiddlewareGroup["/hello"]
	if !ok {
		t.Errorf("expect prefix:/hello's middlewarefunc exists, but not")
	}

	if len(m2) != 1 {

		t.Errorf("expect middlewaregroup len is 1, got %d", len(m2))
	}

	if reflect.TypeOf(m2[0]).String() != reflect.TypeOf(m[0]).String() {
		t.Errorf("expect type  []middlewareFunc, got %s", reflect.TypeOf(m2).String())
	}

	g := app.GetRouteGroup("/hello")
	if reflect.TypeOf(g).String() != reflect.TypeOf(&echo.Group{}).String() {
		t.Errorf("expect group's type is *echo.Group, got %v", g)
	}
}
