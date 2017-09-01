package gopress

import (
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

// App wrapper of echo.Echo and Container
type App struct {
	*echo.Echo

	Logger   *logrus.Logger
	Services *Container
}

// App returns App instance of server
func (s *Server) App() *App {
	return s.app
}
