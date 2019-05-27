package tests

import (
	"fmt"
	"github.com/eddieowens/axon"
)

const TestServiceDependencyInstanceName = "testServiceDependency"

type TestServiceDependency interface {
	DoEvenMoreTestStuff() string
}

type TestServiceDependencyImpl struct {
	DepTwo  DepTwo  `inject:"DepTwo"`
	Int64   int64   `inject:"constantInt64"`
	Int32   int32   `inject:"constantInt32"`
	Float64 float64 `inject:"constantFloat64"`
}

func (t *TestServiceDependencyImpl) DoEvenMoreTestStuff() string {
	return t.DepTwo.CallDepTwo() + " im the dependency! " + fmt.Sprint(t.Int64) + fmt.Sprint(t.Int32) + fmt.Sprint(t.Float64)
}

func TestServiceDependencyFactory(inj axon.Injector, _ axon.Args) axon.Instance {
	t := new(TestServiceDependencyImpl)
	t.DepTwo = inj.GetStructPtr("DepTwo").(DepTwo)
	return axon.StructPtr(t)
}
