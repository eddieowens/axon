package tests

import "axon"

const TestServiceDependencyInstanceName = "testServiceDependency"

type TestServiceDependency interface {
	axon.Instance
	DoEvenMoreTestStuff() string
}

type TestServiceDependencyImpl struct {
	DepTwo DepTwo `inject:"DepTwo"`
}

func (*TestServiceDependencyImpl) GetInstanceName() string {
	return TestServiceDependencyInstanceName
}

func (t *TestServiceDependencyImpl) DoEvenMoreTestStuff() string {
	return t.DepTwo.CallDepTwo() + " im the dependency!"
}

func TestServiceDependencyFactory(_ axon.Args) axon.Instance {
	return new(TestServiceDependencyImpl)
}
