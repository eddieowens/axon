package tests

import (
	. "github.com/eddieowens/axon"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInjector(t *testing.T) {
	asrt := assert.New(t)

	binder := NewBinder(
		NewModule(
			Bind("testService").To().Factory(TestServiceFactory).WithArgs(Args{"arg value", 1, float32(2.0)}),
			Bind("constantString").To().String("the string"),
			Bind("constantInt").To().Int(2),
			Bind("constantInt32").To().Int32(3),
			Bind("constantInt64").To().Int64(4),
			Bind("constantFloat32").To().Float32(5),
			Bind("constantFloat64").To().Float64(6),
			Bind("constantBool").To().Bool(true),
			Bind(TestServiceDependencyInstanceName).To().Factory(TestServiceDependencyFactory).WithoutArgs(),
			Bind(DepTwoInstanceName).To().Instance(StructPtr(new(DepTwoImpl))),
		),
	)

	injector := NewInjector(binder)

	ts := injector.GetStructPtr("testService").(*TestServiceImpl)
	asrt.Equal("dep two! true im the dependency! 436test! the string25", ts.DoTestStuff())
	asrt.NotNil("im the dependency!", injector.GetStructPtr(TestServiceDependencyInstanceName).(TestServiceDependency).DoEvenMoreTestStuff())
	asrt.Equal("arg value", ts.StringField)
	asrt.Equal(1, ts.IntField)
	asrt.Equal(float32(2.0), ts.Float32Field)
	asrt.Equal(uint(0), ts.UIntField)
}

func TestInjectTestServiceDependencyMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.Add(TestServiceDependencyInstanceName, StructPtr(new(MockTestServiceDependency)))

	asrt.Equal("this is a mock!", injector.GetStructPtr(TestServiceDependencyInstanceName).(TestServiceDependency).DoEvenMoreTestStuff())
}

func TestInjectTestServiceMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.AddProvider("testService", NewProvider(TestServiceMockFactory))

	asrt.Equal("I'm a mock provider!", injector.GetStructPtr("testService").(TestService).DoTestStuff())
}

func TestMultipleAdd(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	ts := injector.GetStructPtr("testService").(TestService)
	asrt.Equal("dep two! true im the dependency! 436test! the string25", ts.DoTestStuff())

	injector.Add(TestServiceDependencyInstanceName, StructPtr(new(MockTestServiceDependency)))
	injector.Add(TestServiceDependencyInstanceName, injector.Get(TestServiceDependencyInstanceName))
	injector.GetStructPtr("testService")
	injector.Add(TestServiceDependencyInstanceName, StructPtr(new(MockTestServiceDependency)))
	mock := injector.GetStructPtr("testService").(TestService)

	asrt.Equal("this is a mock!test! the string25", mock.DoTestStuff())
}

func TestAdd(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	injector.GetStructPtr("testService")
	injector.Add(DepTwoInstanceName, StructPtr(new(MockDepTwo)))
	mock := injector.GetStructPtr("testService").(TestService)

	asrt.Equal("mock dep two im the dependency! 436test! the string25", mock.DoTestStuff())
}

func TestStruct(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	injector.Add("something", StructPtr(new(TestServiceImpl)))
	actual := injector.Get("something").GetStructPtr().(*TestServiceImpl).DoTestStuff()

	asrt.Equal("dep two! true im the dependency! 436test! the string25", actual)
}

func TestConstantPrecedence(t *testing.T) {
	asrt := assert.New(t)

	type constStruct struct {
		Int int `inject:"int"`
	}

	binder := NewBinder(
		NewModule(
			Bind("const").To().Factory(func(_ Injector, args Args) Instance {
				return StructPtr(&constStruct{Int: args.Int(0)})
			}).WithArgs(Args{1}),
			Bind("int").To().Int(5),
		),
	)

	injector := NewInjector(binder)

	depTwo := injector.GetStructPtr("const").(*constStruct)
	asrt.Equal(1, depTwo.Int)
}

func createInjector() Injector {

	binder := NewBinder(
		NewModule(
			Bind("constantString").To().String("the string"),
			Bind("constantInt").To().Int(2),
			Bind("constantInt32").To().Int32(3),
			Bind("constantInt64").To().Int64(4),
			Bind("constantFloat32").To().Float32(5),
			Bind("constantFloat64").To().Float64(6),
			Bind("constantBool").To().Bool(true),
			Bind("testService").To().Factory(TestServiceFactory).WithoutArgs(),
			Bind(TestServiceDependencyInstanceName).To().Factory(TestServiceDependencyFactory).WithoutArgs(),
			Bind(DepTwoInstanceName).To().Instance(StructPtr(new(DepTwoImpl))),
		),
	)

	return NewInjector(binder)
}
