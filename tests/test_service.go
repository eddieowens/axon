package tests

import "axon"

type TestService interface {
	axon.Instance
	DoTestStuff() string
}

type TestServiceImpl struct {
	TestServiceDependency TestServiceDependency `inject:"testServiceDependency"`
	SomethingElse         string
	StringField           string
	IntField              int
	Float32Field          float32
	UIntField             uint
}

func TestServiceFactory(args axon.Args) axon.Instance {
	t := TestServiceImpl{
		StringField:  args.String(0),
		IntField:     args.Int(1),
		Float32Field: args.Float32(2),
		UIntField:    args.UInt(3),
	}
	return &t
}

func (TestServiceImpl) GetInstanceName() string {
	return "testService"
}

func (t TestServiceImpl) DoTestStuff() string {
	return t.TestServiceDependency.DoEvenMoreTestStuff() + "test!"
}
