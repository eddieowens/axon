package tests

type MockTestServiceDependency struct {
}

func (*MockTestServiceDependency) DoEvenMoreTestStuff() string {
	return "this is a mock!"
}
