package log

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	contextLoggerKey = "gopress.logging.key"
)

func Extract(ctx echo.Context) *logrus.Entry {
	if logger, ok := ctx.Get(contextLoggerKey).(*logrus.Entry); ok {
		return logger
	}
	return DefaultLogger.WithField("logger", "default")
}

func WithLogger(ctx echo.Context, logger *logrus.Entry) echo.Context {
	ctx.Set(contextLoggerKey, logger)
	return ctx
}
