package gopress

import (
	"testing"
)

type serviceForTest struct {
	name      string
	container *Container
}

func (s *serviceForTest) ServiceName() string {
	return s.name
}

func (s *serviceForTest) RegisterContainer(c *Container) {
	s.container = c
}

func TestContainerRegisterAndGet(t *testing.T) {
	container := NewContainer()

	cases := []struct {
		name    string
		service *serviceForTest
		empty   bool
	}{
		{"a", &serviceForTest{name: "a"}, false},
		{"b", &serviceForTest{name: "b"}, false},
		{"c", &serviceForTest{name: "not c"}, true},
	}

	for _, c := range cases {
		container.Register(c.service)

		actual := container.Get(c.name)
		if c.empty {
			if actual != nil {
				t.Errorf("expect container get value is nil, actual is %#v", actual)
			}
		} else {
			if _, ok := actual.(*serviceForTest); !ok {
				t.Errorf("expect container get value type is *serviceForTest, actual is %#v", actual)
			}
		}
	}
}
