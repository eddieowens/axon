package tests

import (
	. "axon"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInjector(t *testing.T) {
	asrt := assert.New(t)

	binder := []BinderEntry{
		{
			Provider: &Provider{
				Factory:      TestServiceFactory,
				InstanceName: "testService",
				Args:         Args{"arg value", 1, float32(2.0)},
			},
		},
		{
			Provider: &Provider{
				Factory:      TestServiceDependencyFactory,
				InstanceName: TestServiceDependencyInstanceName,
			},
		},
		{
			Instance: new(DepTwoImpl),
		},
	}

	injector := NewInjector(binder)

	asrt.Equal("dep two! im the dependency!test!", injector.GetInstance("testService").(TestService).DoTestStuff())
	asrt.NotNil("im the dependency!", injector.GetInstance(TestServiceDependencyInstanceName).(TestServiceDependency).DoEvenMoreTestStuff())
	asrt.Equal("arg value", injector.GetInstance("testService").(*TestServiceImpl).StringField)
	asrt.Equal(1, injector.GetInstance("testService").(*TestServiceImpl).IntField)
	asrt.Equal(float32(2.0), injector.GetInstance("testService").(*TestServiceImpl).Float32Field)
	asrt.Equal(uint(0), injector.GetInstance("testService").(*TestServiceImpl).UIntField)
}

func TestInjectTestServiceDependencyMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.AddInstance(&MockTestServiceDependency{})

	asrt.Equal("this is a mock!", injector.GetInstance(TestServiceDependencyInstanceName).(TestServiceDependency).DoEvenMoreTestStuff())
}

func TestInjectTestServiceMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.AddProvider("testService", &Provider{Factory: TestServiceMockFactory})

	asrt.Equal("I'm a mock provider!", injector.GetInstance("testService").(TestService).DoTestStuff())
}

func TestMultipleAddInstance(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	ts := injector.GetInstance("testService").(TestService)
	asrt.Equal("dep two! im the dependency!test!", ts.DoTestStuff())

	injector.AddInstance(new(MockTestServiceDependency))
	injector.AddInstance(injector.GetInstance(TestServiceDependencyInstanceName))
	injector.GetInstance("testService")
	injector.AddInstance(new(MockTestServiceDependency))
	mock := injector.GetInstance("testService").(TestService)

	asrt.Equal("this is a mock!test!", mock.DoTestStuff())
}

func TestAddInstance(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	injector.GetInstance("testService")
	injector.AddInstance(new(MockDepTwo))
	mock := injector.GetInstance("testService").(TestService)

	asrt.Equal("mock dep two im the dependency!test!", mock.DoTestStuff())
}

func createInjector() Injector {
	binder := []BinderEntry{
		{
			Provider: &Provider{
				Factory:      TestServiceFactory,
				InstanceName: "testService",
			},
		},
		{
			Provider: &Provider{
				Factory:      TestServiceDependencyFactory,
				InstanceName: TestServiceDependencyInstanceName,
			},
		},
		{
			Instance: new(DepTwoImpl),
		},
	}

	return NewInjector(binder)
}
