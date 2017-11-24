package gopress

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestNewTracingMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	buf := new(bytes.Buffer)
	l := NewLogger()
	l.SetOutput(buf)

	testCases := []struct {
		OpName string
	}{
		{
			OpName: "test",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.OpName, func(t *testing.T) {
			tr := mocktracer.New()
			mw := NewTracingMiddleware(tr, tc.OpName)(func(c Context) error {
				return c.JSON(http.StatusOK, "ok")
			})
			h := NewLoggingMiddleware("test opentracing logger", l)(func(c Context) error {
				return mw(c)
			})
			h(c)

			spans := tr.FinishedSpans()
			if expect, want := len(spans), 1; expect != want {
				t.Fatalf("expected %d spans, expected %d", expect, want)
			}

			if expect, want := spans[0].OperationName, tc.OpName; expect != want {
				t.Fatalf("expect %s operation name, expected %s", expect, want)
			}
		})
	}
}
