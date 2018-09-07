package gopress

import (
	"github.com/fpay/gopress/log"
)

// Logger wraps logrus.Logger
type Logger = log.Logger

var defaultLogger = log.DefaultLogger
