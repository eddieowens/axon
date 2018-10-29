package tests

type MockTestServiceDependency struct {
}

func (*MockTestServiceDependency) GetInstanceName() string {
	return TestServiceDependencyInstanceName
}

func (*MockTestServiceDependency) DoEvenMoreTestStuff() string {
	return "this is a mock!"
}
