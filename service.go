package gopress

import (
	"sync"
)

// Service 服务接口
type Service interface {
	ServiceName() string
	RegisterContainer(c *Container)
}

// Container 服务容器
type Container struct {
	m *sync.Map
}

// NewContainer 创建新的服务容器
func NewContainer() *Container {
	c := &Container{
		m: new(sync.Map),
	}
	return c
}

// Register 注册服务到容器
func (c *Container) Register(svc Service) {
	c.m.Store(svc.ServiceName(), svc)
	svc.RegisterContainer(c)
}

// Get 从容器中获取服务
func (c *Container) Get(name string) Service {
	svc, ok := c.m.Load(name)
	if ok {
		return svc.(Service)
	}
	return nil
}
