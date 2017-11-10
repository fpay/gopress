package gopress

import (
	"net/http"
	"net/http/httptest"
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

	app := &App{
		Logger: defaultLogger,
	}
	var actual Context
	h := NewAppContextMiddleware(app)(func(c Context) error {
		actual = c
		return c.String(http.StatusOK, "test context")
	})
	h(c)

	if c, ok := actual.(*AppContext); !ok {
		t.Errorf("expect context is AppContext, actual is %#v", c)
	}
}
