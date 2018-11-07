package tests

import "axon"

type MockTestService struct {
}

func TestServiceMockFactory(_ axon.Args) axon.Instance {
	return new(MockTestService)
}

func (MockTestService) GetInstanceName() string {
	return "testService"
}

func (MockTestService) DoTestStuff() string {
	return "I'm a mock provider!"
}
