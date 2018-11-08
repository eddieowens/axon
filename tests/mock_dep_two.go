package tests

type MockDepTwo struct {
}

func (*MockDepTwo) GetInstanceName() string {
	return DepTwoInstanceName
}

func (*MockDepTwo) CallDepTwo() string {
	return "mock dep two"
}
