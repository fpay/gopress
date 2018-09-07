package log

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var (
	testLoggingOutput = new(bytes.Buffer)
)

func init() {
	DefaultLogger.SetOutput(testLoggingOutput)
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
