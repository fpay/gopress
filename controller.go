package gopress

// Controller 控制器接口
type Controller interface {

	// RegisterRoutes 注册控制器路由
	RegisterRoutes(app *App)
}
