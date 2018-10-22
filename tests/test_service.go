package tests

import "axon"

type testService interface {
	axon.Instance
	DoTestStuff() string
}

type testServiceImpl struct {
	testServiceDependency testServiceDependency
}

func testServiceFactory(injector axon.Injector, _ ...interface{}) axon.Instance {
	testServiceDependency := injector.GetInstance("testServiceDependency").(testServiceDependency)
	return testServiceImpl{
		testServiceDependency: testServiceDependency,
	}
}

func (testServiceImpl) GetInstanceName() string {
	return "testService"
}

func (testServiceImpl) DoTestStuff() string {
	return "test!"
}
