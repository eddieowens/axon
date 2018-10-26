package tests

import (
	"axon"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInjector(t *testing.T) {
	asrt := assert.New(t)

	binder := []axon.BinderEntry{
		{
			Provider: &axon.Provider{Factory: testServiceFactory},
			Name:     "testService1",
		},
		{
			Instance: &testServiceDependencyImpl{},
		},
	}

	injector := axon.NewInjector(binder)

	asrt.Equal("im the dependency!test!", injector.GetInstance("testService1").(testService).DoTestStuff())
	asrt.NotNil("im the dependency!", injector.GetInstance(testServiceDependencyInstanceName).(testServiceDependency).DoEvenMoreTestStuff())
}

func TestInjectTestServiceDependencyMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.AddInstance(&mockTestServiceDependency{})

	asrt.Equal("this is a mock!", injector.GetInstance(testServiceDependencyInstanceName).(testServiceDependency).DoEvenMoreTestStuff())
}

func TestInjectTestServiceMock(t *testing.T) {
	asrt := assert.New(t)

	injector := createInjector()

	injector.AddProvider("testService1", &axon.Provider{Factory: testServiceMockFactory})

	asrt.Equal("I'm a mock provider!", injector.GetInstance("testService1").(testService).DoTestStuff())
}

func TestMultipleAddInstance(t *testing.T) {
	asrt := assert.New(t)
	injector := createInjector()

	ts := injector.GetInstance("testService1").(testService)
	asrt.Equal("im the dependency!test!", ts.DoTestStuff())

	injector.AddInstance(new(mockTestServiceDependency))
	mock := injector.GetInstance("testService1").(testService)

	asrt.Equal("this is a mock!test!", mock.DoTestStuff())
}

func createInjector() axon.Injector {
	binder := []axon.BinderEntry{
		{
			Provider: &axon.Provider{Factory: testServiceFactory},
			Name:     "testService1",
		},
		{
			Instance: &testServiceDependencyImpl{},
		},
	}

	return axon.NewInjector(binder)
}
