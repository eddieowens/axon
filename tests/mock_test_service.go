package tests

import (
	"github.com/eddieowens/axon"
)

type MockTestService struct {
}

func TestServiceMockFactory(_ axon.Injector, _ axon.Args) axon.Instance {
	return axon.StructPtr(new(MockTestService))
}

func (MockTestService) DoTestStuff() string {
	return "I'm a mock provider!"
}
