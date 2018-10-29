package tests

import "axon"

const TestServiceDependencyInstanceName = "testServiceDependency"

type TestServiceDependency interface {
	axon.Instance
	DoEvenMoreTestStuff() string
}

type TestServiceDependencyImpl struct {
}

func (*TestServiceDependencyImpl) GetInstanceName() string {
	return TestServiceDependencyInstanceName
}

func (*TestServiceDependencyImpl) DoEvenMoreTestStuff() string {
	return "im the dependency!"
}
