package gopress

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var (
	testLoggingOutput = new(bytes.Buffer)
)

func init() {
	defaultLogger.SetOutput(testLoggingOutput)
}

func TestLoggerSetOutput(t *testing.T) {
	l := &Logger{logrus.StandardLogger()}

	cases := []io.Writer{
		os.Stdout,
		os.Stderr,
		new(bytes.Buffer),
	}
	for _, w := range cases {
		l.SetOutput(w)
		actual := l.Logger.Out
		if actual != w {
			t.Errorf("expect logger output is %#v, actual is %#v", w, actual)
		}
	}
}

func TestLoggerSetFormatter(t *testing.T) {
	l := &Logger{logrus.StandardLogger()}

	cases := []logrus.Formatter{
		&logrus.JSONFormatter{},
		&logrus.TextFormatter{},
	}

	for _, f := range cases {
		l.SetFormatter(f)
		actual := l.Logger.Formatter
		if actual != f {
			t.Errorf("expect logger formatter is %#v, actual is %#v", f, actual)
		}
	}
}

func TestNewLogger(t *testing.T) {
	l := NewLogger()
	if l.Logger.Level != defaultLoggingLevel {
		t.Errorf("expect logging level is %d, actual is %d", defaultLoggingLevel, l.Logger.Level)
	}
	if l.Logger.Out != defaultLoggingOutput {
		t.Errorf("expect logging output is %#v, actual is %#v", defaultLoggingOutput, l.Logger.Out)
	}
	if l.Logger.Formatter != defaultLoggingFormatter {
		t.Errorf("expect logging formatter is %#v, actual is %#v", defaultLoggingFormatter, l.Logger.Formatter)
	}
}

func TestNewLoggingMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	var l *Logger
	var h HandlerFunc
	var buf *bytes.Buffer

	// test with global logger
	testLoggingOutput.Reset()
	h = NewLoggingMiddleware("default logger", nil)(func(c Context) error {
		return c.String(http.StatusOK, "test")
	})
	h(c)

	if testLoggingOutput.Len() == 0 {
		t.Errorf("expect test logging output in global buffer not empty")
	}
	testLoggingOutput.Reset()

	// test with custom logger
	buf = new(bytes.Buffer)
	l = NewLogger()
	l.SetOutput(buf)
	h = NewLoggingMiddleware("test logger", l)(func(c Context) error {
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
	app := &App{Logger: NewLogger()}
	app.Logger.SetOutput(buf)
	h = NewLoggingMiddleware("app logger", nil)(func(c Context) error {
		return c.String(http.StatusOK, "test")
	})
	h(&AppContext{c, app})

	if buf.Len() == 0 {
		t.Errorf("expect test logging output in app buffer not empty")
	}
}
