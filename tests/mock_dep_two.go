package tests

type MockDepTwo struct {
}

func (*MockDepTwo) CallDepTwo() string {
	return "mock dep two"
}
