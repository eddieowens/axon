package tests

import "axon"

type TestService interface {
	axon.Instance
	DoTestStuff() string
}

type TestServiceImpl struct {
	TestServiceDependency TestServiceDependency `inject:"testServiceDependency"`
}

func TestServiceFactory() axon.Instance {
    return new(TestServiceImpl)
}

func (TestServiceImpl) GetInstanceName() string {
	return "testService"
}

func (t TestServiceImpl) DoTestStuff() string {
	return t.TestServiceDependency.DoEvenMoreTestStuff() + "test!"
}
