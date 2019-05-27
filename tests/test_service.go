package tests

import (
	"fmt"
	"github.com/eddieowens/axon"
)

type TestService interface {
	DoTestStuff() string
}

type TestServiceImpl struct {
	TestServiceDependency TestServiceDependency `inject:"testServiceDependency"`
	SomethingElse         string                `inject:"constantString"`
	StringField           string
	InjectedIntField      int `inject:"constantInt"`
	IntField              int
	Float32Field          float32
	InjectedFloat32Field  float32 `inject:"constantFloat32"`
	UIntField             uint
}

func TestServiceFactory(_ axon.Injector, args axon.Args) axon.Instance {
	t := TestServiceImpl{
		StringField:  args.String(0),
		IntField:     args.Int(1),
		Float32Field: args.Float32(2),
		UIntField:    args.UInt(3),
	}
	return axon.StructPtr(&t)
}

func (t TestServiceImpl) DoTestStuff() string {
	return t.TestServiceDependency.DoEvenMoreTestStuff() + "test! " + t.SomethingElse + fmt.Sprint(t.InjectedIntField) + fmt.Sprint(t.InjectedFloat32Field)
}
