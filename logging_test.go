package gopress

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var (
	testLoggingOutput = new(bytes.Buffer)
)

func init() {
	defaultLogger.SetOutput(testLoggingOutput)
}

func TestLoggerSetOutput(t *testing.T) {
	l := &Logger{logrus.StandardLogger(), defaultLoggingLevel}

	cases := []io.Writer{
		os.Stdout,
		os.Stderr,
		new(bytes.Buffer),
	}
	for _, w := range cases {
		l.SetOutput(w)
		actual := l.Output()
		if actual != w {
			t.Errorf("expect logger output is %#v, actual is %#v", w, actual)
		}
		if l.Logger.Out != w {
			t.Errorf("expect underlying output is %#v, actual is %#v", w, l.Logger.Out)
		}
	}
}

func TestLoggerSetFormatter(t *testing.T) {
	l := &Logger{logrus.StandardLogger(), defaultLoggingLevel}

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
	if l.level != defaultLoggingLevel {
		t.Errorf("expect logging level is %d, actual is %d", defaultLoggingLevel, l.Logger.Level)
	}
	if l.Logger.Out != defaultLoggingOutput {
		t.Errorf("expect logging output is %#v, actual is %#v", defaultLoggingOutput, l.Logger.Out)
	}
	if l.Logger.Formatter != defaultLoggingFormatter {
		t.Errorf("expect logging formatter is %#v, actual is %#v", defaultLoggingFormatter, l.Logger.Formatter)
	}
}

func TestLoggingPrefix(t *testing.T) {
	l := NewLogger()

	expect := ""
	cases := []string{"a", "b", "c", "gopress", "echo"}
	for _, c := range cases {
		l.SetPrefix(c)
		actual := l.Prefix()
		if actual != expect {
			t.Errorf("expect prefix is %v, actual is %v", expect, actual)
		}
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

	// test with handler error
	buf.Reset()
	e.Logger = app.Logger
	h = NewLoggingMiddleware("handler error", nil)(func(c Context) error {
		return errors.New("test error")
	})
	h(&AppContext{c, app})

	if buf.Len() == 0 {
		t.Errorf("expect test logging output in app buffer not empty")
	}

	if !bytes.Contains(buf.Bytes(), []byte(`"level":"error"`)) {
		t.Errorf("expect test loggint contains level error")
	}
}

func TestLoggerLevel(t *testing.T) {
	l := NewLogger()

	cases := []log.Lvl{
		log.DEBUG,
		log.INFO,
		log.WARN,
		log.ERROR,
		log.OFF,
	}

	for _, v := range cases {
		l.SetLevel(v)
		if l.Level() != v {
			t.Errorf("expect logging level is %v, actual is %v", v, l.level)
		}

		logrusLevel := loggingLevelMapping[v]
		if l.Logger.Level != logrusLevel {
			t.Errorf("expect underlying logrus.Logger's level is %v, actual is %v", logrusLevel, l.Logger.Level)
		}
	}
}
