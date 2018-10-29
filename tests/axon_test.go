package tests

import (
    "github.com/stretchr/testify/assert"
    "testing"
    . "axon"
)

func TestNewInjector(t *testing.T) {
    asrt := assert.New(t)

    binder := []BinderEntry{
        {
            Provider: &Provider{
                Factory:      TestServiceFactory,
                InstanceName: "testService",
            },
        },
        {
            Instance: &TestServiceDependencyImpl{},
        },
    }

    injector := NewInjector(binder)

    asrt.Equal("im the dependency!test!", injector.GetInstance("testService").(TestService).DoTestStuff())
    asrt.NotNil("im the dependency!", injector.GetInstance(TestServiceDependencyInstanceName).(TestServiceDependency).DoEvenMoreTestStuff())
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
    asrt.Equal("im the dependency!test!", ts.DoTestStuff())

    injector.AddInstance(new(MockTestServiceDependency))
    injector.AddInstance(injector.GetInstance(TestServiceDependencyInstanceName))
    injector.GetInstance("testService")
    injector.AddInstance(new(MockTestServiceDependency))
    mock := injector.GetInstance("testService").(TestService)

    asrt.Equal("this is a mock!test!", mock.DoTestStuff())
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
            Instance: &TestServiceDependencyImpl{},
        },
    }

    return NewInjector(binder)
}
