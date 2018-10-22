# axon
A simple, lightweight, lazy-loaded, reflectionless, and concurrent DI (really just a singleton management) library.

## Install
```bash
go get github.com/eddieowens/axon
```

## Usage
### Basic
Define some interface with the `axon.Instance` interface embedded
```go
type TestService interface {
    axon.Instance
    DoTestStuff() string
}
```

Implement the interface
```go
const TestServiceInstance = "testService"

type TestServiceImpl struct {
}

func (TestServiceImpl) GetInstanceName() string {
    return TestServiceInstance
}

func (TestServiceImpl) DoTestStuff() string {
    return "test!"
}
```
The `axon.Instance` interface requires a single method which defines your service's referent for your 
`axon.Instance` in the `axon.Injector`

Instantiate the `Injector` for your app
```go
package main

var Injector axon.Injector

func main() {
    Injector = CreateInjector()
}

func CreateInjector() axon.Injector {
    binder := []axon.BinderEntry{
        {Instance: testServiceImpl{}},
    }
    
    return axon.NewInjector(binder)
}
```
The `axon.Injector` is a manager for singletons created as defined within your `axon.Binder` which is
just an alias for a slice of `axon.BinderEntry`. This binder is customizable. Check out the docs to learn
more about customizing the binder your injector manages

Now get your newly created service!
```go
func SomethingToDoInMyApp() {
    testService := Instance.GetInstance(TestServiceInstance).(TestService)
    fmt.Println(testService.DoTestStuff()) // Prints "test!"
}
```

### In tests
In order to mock the `TestService` within a test, define your mock
```go
type TestServiceMock struct {
}

func (TestServiceMock) GetInstanceName() string {
    return TestServiceInstance
}

func (TestServiceMock) DoTestStuff() string {
    return "I'm a mock!"
}
```

Create the injector and override the `TestServiceImpl` instance with your mock
```go
func TestSomethingThatUsesTestService(t *testing.T) {
    injector := CreateInjector() // don't worry, this call is incredibly light even for very large binders
    injector.AddInstance(TestServiceMock{})
    ...
}
```
Now wherever `injector.GetInstance("testService")` is called, the mock will be returned.

## Why should I use this?
Many dependency injection frameworks out there use quite a bit of heavy and brittle reflection to achieve 
the desired state. Axon is simply a singleton manager which eliminates a lot of boiler plate and error prone
code (especially around [concurrency](http://marcio.io/2015/07/singleton-pattern-in-go/)). It uses no reflection,
lazy-loads dependencies, and is safe to use in multiple goroutines. It's also very easy to override 
instances within the injector with a mock for easy testing.

In short, it makes managing a GoLang codebase a bit easier.

## [Docs](https://godoc.org/github.com/eddieowens/axon)

## License
[MIT](https://github.com/eddieowens/axon/blob/master/LICENSE)