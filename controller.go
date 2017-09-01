package gopress

// Controller 控制器接口
type Controller interface {

	// RegisterRoutes 注册控制器路由
	RegisterRoutes(app *App)
}

// RegisterController 注册控制器
func (s *Server) RegisterController(c Controller) {
	c.RegisterRoutes(s.app)
}

// RegisterControllers 注册控制器列表
func (s *Server) RegisterControllers(cs ...Controller) {
	for _, c := range cs {
		s.RegisterController(c)
	}
}
