# axon
A simple, lightweight, and lazy-loaded DI (really just a singleton management) library.

## Install
```bash
go get github.com/eddieowens/axon
```

## Usage
### Basic
```go
package main

import (
    "fmt"
    "github.com/eddieowens/axon"
)

func main() {
    binder := axon.NewBinder(axon.NewModule(
        axon.Bind("AnswerToTheUltimateQuestion").To().Int(42),
    ))
    
    injector := axon.NewInjector(binder)
    
    fmt.Print(injector.GetInt("AnswerToTheUltimateQuestion")) // Prints 42
}
```
In the above, I created a new `Binder` by passing in a `Module` which stores my int value (42). A `Binder` is a series
of `Modules` and these `Modules` allow you to define what is stored in the `Injector` at runtime. You define a `Module` by
using different `Bindings`. In this case, we bound the `key` _AnswerToTheUltimateQuestion_ to the int value 42. Now on all subsequent
`injector.GetInt("AnswerToTheUltimateQuestion")` calls, 42 will be returned.
### Injecting dependencies
Now the above isn't very interesting on its own but what if you wanted to pass in the `AnswerToTheUltimateQuestion` value to everything that
depended on it? It would be tedious to create a global variable and pass that along wherever you needed it. Instead,
you can `inject` it as a dependency.
```go
package main

import (
    "fmt"
    "github.com/eddieowens/axon"
)

func main() {
    type MyStruct struct {
        IntField int `inject:"AnswerToTheUltimateQuestion"`
    }
    
    binder := axon.NewBinder(axon.NewModule(
        axon.Bind("MyStruct").To().Instance(axon.StructPtr(new(MyStruct))),
        axon.Bind("AnswerToTheUltimateQuestion").To().Int(42),
    ))
    
    injector := axon.NewInjector(binder)
    
    fmt.Print(injector.GetStructPtr("MyStruct").(*MyStruct).IntField) // Prints 42
}

```
Now, you have a struct called `MyStruct` which is bound to the `key` _MyStruct_ as an `Instance`. An `Instance` is the
wrapper that `axon` uses around everything it manages. It holds metadata and allows the `injector` to interact with 
everything you've defined in your `Binder` efficiently and safely. All interactions between your raw data and the 
`injector` will be done via an `Instance`.

You may have also noticed the `inject` tag. Well this is how `axon` delivers your dependencies to your struct at
runtime and allows you to not have to worry about creating it yourself. The `inject` tag takes a single string value
which is the `key` you defined within your `Module` and will automatically pass the value in on a call to 
`injector.GetStructPtr("MyStruct")`.

### Utilizing interfaces
In the above example, you have to write quite a bit of code to do something that could be done pretty succinctly
in just raw go. Admittedly, `axon` does not shine when injecting simple constants (although it does make managing them
far easier). Where `axon` really shines is injecting deep dependency trees and using interfaces.
```go
package main

import (
    "fmt"
    "github.com/eddieowens/axon"
)

type Starter interface {
    Start()
}

type Car struct {
    Engine Starter `inject:"Engine"`
}

func (c *Car) Start() {
    fmt.Println("Starting the Car!")
    c.Engine.Start()
}

type Engine struct {
    FuelInjector Starter `inject:"FuelInjector"`
}

func (e *Engine) Start() {
    fmt.Println("Starting the Engine!")
    e.FuelInjector.Start()
}

type FuelInjector struct {
}

func (*FuelInjector) Start() {
    fmt.Println("Starting the FuelInjector!")
}

func CarFactory(_ axon.Args) axon.Instance {
    fmt.Println("Hey, a new Car is being made!")
    return axon.StructPtr(new(Car))
}

func main() {
    binder := axon.NewBinder(axon.NewModule(
        axon.Bind("Car").To().Factory(CarFactory).WithoutArgs(),
        axon.Bind("Engine").To().Instance(axon.StructPtr(new(Engine))),
        axon.Bind("FuelInjector").To().Instance(axon.StructPtr(new(FuelInjector))),
    ))
    
    injector := axon.NewInjector(binder)
    
    // Prints:
    // Hey, a new Car is being made!
    // Starting the Car!
    // Starting the Engine!
    // Starting the FuelInjector!
    injector.GetStructPtr("Car").(Starter).Start()
}

```
Here we have a `Car` which depends on an `Engine` which depends on a `FuelInjector`. Things are getting a bit messy
and managing all of these structs is becoming tedious and cumbersome. Rather than managing everything ourselves, `axon`
can manage these dependencies for us. Now whenever you call the `Start()` function on the `Car`, all of the required 
dependencies will be automatically added and this will be consistent throughout your codebase.

### Factories
The above also introduces the use of a `Factory` to create an `Instance`. You use a `Factory` when you want the
construction of your `Instance` to be managed by `axon`. For instance, let's say you require some parameters for
your `Car` that aren't provided until runtime and aren't managed by `axon`.
```go
...
type Car struct {
    LockCode string
}

func CarFactory(args axon.Args) axon.Instance {
    return &Car{
        LockCode: args.String(0),
    }
}
...
binder := axon.NewBinder(axon.NewModule(
    axon.Bind("Car").To().Factory(CarFactory).WithArgs(axon.Args{os.Getenv("CAR_LOCK_CODE")}),
    ...
))

injector := axon.NewInjector(binder)

car := injector.GetStructPtr("Car").(*Car)
fmt.Println(car.LockCode) // Prints the value of env var CAR_LOCK_CODE
```
FYI, if the `Arg` passed into your `Instance` is overriding a field that is tagged with `inject`, the arg will
always take precedence and will not be overwritten. 
### Testing
What I think to be the most useful functionality that `axon` affords you is the ability to very easily and precisely
test your code. Let's say you create a test for your `Engine`. You don't want to also
test the functionality of the `FuelInjector` so you make a mock of the `Starter` interface.
```go
import (
    "fmt"
    "github.com/stretchr/testify/mock"
    "testing"
)

type MockFuelInjector struct {
    mock.Mock
}

func (m *MockFuelInjector) Start() {
    fmt.Println("I'm a mock FuelInjector!")
}
```
Then you add it to the `Injector` within your test.
```go
func TestEngine(t *testing.T) {
    injector := axon.NewInjector(binder)
	
    injector.Add("FuelInjector", axon.StructPtr(new(MockFuelInjector)))
    
    // Prints:
    // Starting the Engine!
    // I'm a mock FuelInjector!
    injector.GetStructPtr("Engine").(Starter).Start()
}
```
The `injector.Add()` method will replace the `Instance` held by the `key` with the mocked `FuelInjector` allowing you
to unit test your code efficiently and exactly.

To note, the `Injector` keeps track of every `Instance`'s dependencies and will clear the dependencies of a particular
`key` when the `Add()` method is called. This means all subsequent `injector.Get()` method calls will return the
correct `Instance`. For example, if you call `injector.Get("Car")` in the above test, it will store an `Engine` which
will store a `MockFuelInjector`.
## [Docs](https://godoc.org/github.com/eddieowens/axon)

## License
[MIT](https://github.com/eddieowens/axon/blob/master/LICENSE)