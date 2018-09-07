package middlewares

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fpay/gopress"
	"github.com/fpay/gopress/log"
	"github.com/labstack/echo"
)

var (
	testLoggingOutput = new(bytes.Buffer)
)

func TestNewLoggingMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	var l *log.Logger
	var h echo.HandlerFunc
	var buf *bytes.Buffer

	// test with global logger
	testLoggingOutput.Reset()
	h = NewLoggingMiddleware(LoggingMiddlewareConfig{
		Name: "default logger",
	})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
	h(c)

	if testLoggingOutput.Len() == 0 {
		t.Errorf("expect test logging output in global buffer not empty")
	}
	testLoggingOutput.Reset()

	// test with custom logger
	buf = new(bytes.Buffer)
	l = log.NewLogger()
	l.SetOutput(buf)
	h = NewLoggingMiddleware(LoggingMiddlewareConfig{
		Name:   "test logger",
		Logger: l,
	})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
	h(c)

	if buf.Len() == 0 {
		t.Errorf("expect test logging output in function buffer not empty")
	}
	if testLoggingOutput.Len() > 0 {
		t.Errorf("expect test logging output in global buffer empty")
	}

	// test with app logger
	buf = new(bytes.Buffer)
	app := &gopress.App{Logger: log.NewLogger()}
	app.Logger.SetOutput(buf)
	h = NewLoggingMiddleware(LoggingMiddlewareConfig{
		Name: "app logger",
	})(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
	h(gopress.NewAppContext(c, app, app.Logger.WithFields(nil)))

	if buf.Len() == 0 {
		t.Errorf("expect test logging output in app buffer not empty")
	}

	// test with handler error
	buf.Reset()
	e.Logger = app.Logger
	h = NewLoggingMiddleware(LoggingMiddlewareConfig{
		Name: "handler error",
	})(func(c echo.Context) error {
		return errors.New("test error")
	})
	h(gopress.NewAppContext(c, app, app.Logger.WithFields(nil)))

	if buf.Len() == 0 {
		t.Errorf("expect test logging output in app buffer not empty")
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"error":"test error"`)) {
		t.Errorf("expect test logging contains (%s)", `"error":"test error"`)
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"level":"error"`)) {
		t.Errorf("expect test loggint contains level error")
	}
}
