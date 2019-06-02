package tests

import (
	"github.com/eddieowens/axon"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InjectorTest struct {
	suite.Suite
}

func (i *InjectorTest) SetupTest() {
}

type Car struct {
	Engine *Engine `inject:"Engine"`
}

type Engine struct {
	FuelInjector *FuelInjector `inject:"FuelInjector"`
}

type FuelInjector struct {
}

func (i *InjectorTest) TestPropagatedInjection() {
	// -- Given
	//
	engine := &Engine{
		FuelInjector: &FuelInjector{},
	}
	binder := axon.NewBinder(axon.NewPackage(
		axon.Bind("Car").To().StructPtr(new(Car)),
		axon.Bind("Engine").To().StructPtr(engine),
	))

	inj := axon.NewInjector(binder)

	// -- When
	//
	car := inj.GetStructPtr("Car").(*Car)

	// -- Then
	//
	i.NotNil(car)
}

func TestInjectorTest(t *testing.T) {
	suite.Run(t, new(InjectorTest))
}
