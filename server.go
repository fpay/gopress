package gopress

import (
	"context"
	"fmt"

	"github.com/fpay/gopress/log"
	"github.com/labstack/echo/v4"
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
	Logger *Logger

	app    *App
	listen string
}

// StaticOptions 静态文件配置
type StaticOptions struct {
	Path string `yaml:"path" mapstructure:"path"`
	Root string `yaml:"root" mapstructure:"root"`
}

// ServerOptions 服务器配置
type ServerOptions struct {
	Host   string        `yaml:"host" mapstructure:"path"`
	Port   int           `yaml:"port" mapstructure:"port"`
	Views  string        `yaml:"views" mapstructure:"views"`
	Static StaticOptions `yaml:"static" mapstructure:"static"`
}

// NewServer 创建HTTP服务器
func NewServer(options ServerOptions) *Server {
	tplRoot := options.Views
	if len(tplRoot) == 0 {
		tplRoot = defaultViewsRoot
	}

	staticPath := options.Static.Path
	if len(staticPath) == 0 {
		staticPath = defaultStaticPath
	}
	staticRoot := options.Static.Root
	if len(staticRoot) == 0 {
		staticRoot = defaultStaticRoot
	}

	port := options.Port
	if options.Port == 0 {
		port = defaultPort
	}

	logger := log.NewLogger()

	e := echo.New()
	e.Logger = logger

	app := &App{
		Echo:     e,
		Logger:   logger,
		Services: NewContainer(),
	}

	app.Renderer = NewTemplateRenderer(tplRoot)
	app.Static(staticPath, staticRoot)
	app.Use(NewAppContextMiddleware(app))

	return &Server{
		Logger: logger,
		app:    app,
		listen: fmt.Sprintf("%s:%d", options.Host, port),
	}
}

// App returns App instance of server
func (s *Server) App() *App {
	return s.app
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	return s.app.Start(s.listen)
}

// StartTLS 启动HTTPS服务器
func (s *Server) StartTLS(cert, key string) error {
	return s.app.StartTLS(s.listen, cert, key)
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.Shutdown(ctx)
}

// RegisterControllers 注册控制器列表
func (s *Server) RegisterControllers(cs ...Controller) {
	for _, c := range cs {
		c.RegisterRoutes(s.app)
	}
}

// RegisterGlobalMiddlewares 注册全局中间件
func (s *Server) RegisterGlobalMiddlewares(middlewares ...MiddlewareFunc) {
	s.app.Use(middlewares...)
}

// RegisterServices 注册服务到server app的服务容器
func (s *Server) RegisterServices(services ...Service) {
	for _, svc := range services {
		s.app.Services.Register(svc)
	}
}
