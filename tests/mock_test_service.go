package tests

import "axon"

type mockTestService struct {
}

func testServiceMockFactory(injector axon.Injector, _ ...interface{}) axon.Instance {
	return mockTestService{}
}

func (mockTestService) GetInstanceName() string {
	return "testService"
}

func (mockTestService) DoTestStuff() string {
	return "I'm a mock provider!"
}
