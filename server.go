package gopress

import (
	"fmt"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

const (
	defaultViewsRoot  = "./views"
	defaultStaticPath = "/static"
	defaultStaticRoot = "./static"
	defaultPort       = 3000
)

// Context is alias of echo.Context
type Context = echo.Context

// MiddlewareFunc is alias of echo.MiddlewareFunc
type MiddlewareFunc = echo.MiddlewareFunc

// HandlerFunc is alias of echo.HandlerFunc
type HandlerFunc = echo.HandlerFunc

// Server HTTP服务器
type Server struct {
	app    *App
	listen string
}

// ServerOptions 服务器配置
type ServerOptions struct {
	Host   string `yaml:"host" mapstructure:"path"`
	Port   int    `yaml:"port" mapstructure:"port"`
	Views  string `yaml:"views" mapstructure:"views"`
	Static struct {
		Path string `yaml:"path" mapstructure:"path"`
		Root string `yaml:"root" mapstructure:"root"`
	} `yaml:"static" mapstructure:"static"`
}

// NewServer 创建HTTP服务器
func NewServer(options ServerOptions) *Server {
	app := &App{
		Echo:     echo.New(),
		Logger:   log.StandardLogger(),
		Services: NewContainer(),
	}

	tplRoot := options.Views
	if len(tplRoot) == 0 {
		tplRoot = defaultViewsRoot
	}
	app.Renderer = NewTemplateRenderer(tplRoot)

	staticPath := options.Static.Path
	if len(staticPath) == 0 {
		staticPath = defaultStaticPath
	}
	staticRoot := options.Static.Root
	if len(staticRoot) == 0 {
		staticRoot = defaultStaticRoot
	}
	app.Static(staticPath, staticRoot)

	port := options.Port
	if options.Port == 0 {
		port = defaultPort
	}

	return &Server{
		app:    app,
		listen: fmt.Sprintf("%s:%d", options.Host, port),
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	return s.app.Start(s.listen)
}

// StartTLS 启动HTTPS服务器
func (s *Server) StartTLS(cert, key string) error {
	return s.app.StartTLS(s.listen, cert, key)
}

// RegisterGlobalMiddlewares 注册全局中间件
func (s *Server) RegisterGlobalMiddlewares(middlewares ...MiddlewareFunc) {
	s.app.Use(middlewares...)
}
