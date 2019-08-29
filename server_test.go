package gopress

import (
	"fmt"
	"testing"
)

type controllerForTest struct {
	app *App
}

func (c *controllerForTest) RegisterRoutes(app *App) {
	c.app = app
}

func TestNewServer(t *testing.T) {
	var s *Server

	// defaults
	s = NewServer(ServerOptions{})

	expectListen := fmt.Sprintf(":%d", defaultPort)
	if s.listen != expectListen {
		t.Errorf("expect server listen %s, actual is %s", expectListen, s.listen)
	}

	viewsRoot := s.app.Renderer.(*TemplateRenderer).root
	if viewsRoot != defaultViewsRoot {
		t.Errorf("expect server views root is %s, actual is %s", defaultViewsRoot, viewsRoot)
	}

	// customs
	s = NewServer(ServerOptions{
		Host:   "127.0.0.1",
		Port:   8000,
		Views:  "./templates",
		Static: StaticOptions{"/public", "./public"},
	})

	expectListen = "127.0.0.1:8000"
	if s.listen != expectListen {
		t.Errorf("expect server listen %s, actual is %s", expectListen, s.listen)
	}

	if viewsRoot := s.app.Renderer.(*TemplateRenderer).root; viewsRoot != "./templates" {
		t.Errorf("expect server views root is %s, actual is %s", "./templates", viewsRoot)
	}
}

func TestRegisterControllers(t *testing.T) {
	s := NewServer(ServerOptions{})
	c := &controllerForTest{}

	s.RegisterControllers(c)
	if c.app != s.App() {
		t.Errorf("expect server app passed to controller")
	}
}
