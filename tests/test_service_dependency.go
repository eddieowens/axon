package tests

import "axon"

const testServiceDependencyInstanceName = "testServiceDependency"

type testServiceDependency interface {
	axon.Instance
	DoEvenMoreTestStuff() string
}

type testServiceDependencyImpl struct {
}

func (*testServiceDependencyImpl) GetInstanceName() string {
	return testServiceDependencyInstanceName
}

func (*testServiceDependencyImpl) DoEvenMoreTestStuff() string {
	return "im the dependency!"
}
