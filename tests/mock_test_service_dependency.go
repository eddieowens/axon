package tests

type mockTestServiceDependency struct {
}

func (*mockTestServiceDependency) GetInstanceName() string {
	return testServiceDependencyInstanceName
}

func (*mockTestServiceDependency) DoEvenMoreTestStuff() string {
	return "this is a mock!"
}
