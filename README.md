# axon
A simple, lightweight, lazy-loaded, and concurrent DI (really just a singleton management) library.

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

func (*TestServiceImpl) GetInstanceName() string {
    return TestServiceInstance
}

func (*TestServiceImpl) DoTestStuff() string {
    return "test!"
}
```
The `Instance` interface requires a single method which defines your service's referent for your 
`Instance` in the `Injector`

Instantiate the `Injector` for your app
```go
package main

import "github.com/eddieowens/axon"

var Injector axon.Injector

func main() {
    Injector = CreateInjector()
}

func CreateInjector() axon.Injector {
    binder := []BinderEntry{
        {Instance: &TestServiceImpl{}},
    }
    
    return axon.NewInjector(binder)
}
```
The `Injector` is a manager for singletons created as defined within your `Binder` which is
just an alias for a slice of `BinderEntry`. This binder is customizable. Check out the docs to learn
more about customizing the binder your injector manages

Now get your newly created service!
```go
func SomethingToDoInMyApp() {
    testService := Instance.GetInstance(TestServiceInstance).(TestService)
    fmt.Println(testService.DoTestStuff()) // Prints "test!"
}
```
### Injecting dependencies
Let's say you have another service which requires the previously created `TestService` so you create a struct
that looks something like this
```go
type HigherLevelService struct {
    TestService TestService
}
```
This is perfectly valid Go and with a single dependency this isn't too bad. Although, you still have to manage the
singletons of `HigherLevelService` and `TestServiceImpl` unless you want to create them each time which is 
inefficient. And what if you need more than one or two dependencies? This can get messy and bug prone fast.

Rather than managing the entire dependency graph yourself, `axon` can do that for you
```go
type HigherLevelService interface {
    axon.Instance
    DoSomething() string
}

type HigherLevelServiceImpl struct {
    TestService TestService `inject:"testService"` // The instance name of TestService
}

func (*HigherLevelServiceImpl) GetInstanceName() string {
    return "higherLevelService"
} 

type (*HigherLevelServiceImpl) DoSomething() string {
    return "I'm doing something!"
} 

// The factory that is used within a Provider in a Binder. This factory will be called 
// when injector.GetInstance("higherLevelService") is called and a ptr to a HigherLevelServiceImpl
// will be returned
func HigherLevelServiceFactory(args axon.Args) axon.Instance {
    return new(HigherLevelServiceImpl)
}
```
Now add both the `HigherLevelService` and the `TestService` to your `Binder`
```go
package main

import "github.com/eddieowens/axon"

var Injector axon.Injector

func main() {
    Injector = CreateInjector()
}

func CreateInjector() axon.Injector {
    binder := []BinderEntry{
        {
            Provider: &Provider{
                Factory:      HigherLevelServiceFactory,
                InstanceName: "higherLevelService", // The instance name must be given when using a Provider in a BinderEntry
            },
        },
        {
            Instance: &TestServiceImpl{},
        },
    }
    
    return axon.NewInjector(binder)
}
```
Now whenever you call 
```go
hls := Injector.GetInstance("higherLevelService").(HigherLevelService)
```
The `HigherLevelServiceImpl` will be returned with a `TestServiceImpl` automatically injected via the `inject` tag
into its `TestService` field. This allows for a nice decoupling between your structs and their dependencies.

### In tests
In order to mock the `TestService` within a test, define your mock
```go
type TestServiceMock struct {
}

func (*TestServiceMock) GetInstanceName() string {
    return TestServiceInstance
}

func (*TestServiceMock) DoTestStuff() string {
    return "I'm a mock!"
}
```
Create the injector and override the `TestServiceImpl` instance with your mock
```go
func TestHigherLevelServiceImpl(t *testing.T) {
    injector := CreateInjector() // don't worry, this call is incredibly light even for very large binders
    injector.AddInstance(&TestServiceMock{})
}
```
Now wherever `injector.GetInstance("testService")` is called, the mock will be returned. This is useful for
testing code that depends on `TestService` like `HigherLevelService` from the previous section
```go
func TestHigherLevelServiceImpl(t *testing.T) {
    injector := CreateInjector()
    injector.AddInstance(&TestServiceMock{})
    
    // Grab the HigherLevelServiceImpl to test it with a mock TestService
    hls := injector.GetInstance("higherLevelService").(*HigherLevelServiceImpl)
    
    // Test the HigherLevelServiceImpl
    ...
}
```

## Why should I use this?
Many dependency injection frameworks out there use quite a bit of heavy and brittle reflection to achieve 
the desired state. Axon is simply a singleton manager which eliminates a lot of boiler plate and error prone
code (especially around [concurrency](http://marcio.io/2015/07/singleton-pattern-in-go/)). It uses very little 
reflection, lazy-loads dependencies, and is safe to use in multiple goroutines. It's also very easy to override 
instances within the injector with a mock for easy testing.

In short, it makes managing a GoLang codebase a bit easier.

## [Docs](https://godoc.org/github.com/eddieowens/axon)

## License
[MIT](https://github.com/eddieowens/axon/blob/master/LICENSE)