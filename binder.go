package axon

import (
	"sync"
)

// Function that constructs an Instance. The args field are provided via the Provider listed
// within a Binding. The Injector provided is the state of the Injector at the time of calling
// the factory. This should be used carefully as a user can end up in an infinite loop
type Factory func(inj Injector, args Args) Instance

// Args that are passed into a Factory upon creation (call to injector.Get*(...)). These Args will take precedence over
// any and all other entities managed by axon. For instance, if you have a field that is tagged with inject, and is
// also instantiated via Args passed in a Factory, the value provided by the Arg will always remain.
//
// Add Args through a Binding like so.
//   axon.NewPackage(axon.Bind("MyService").To().Factory(MyServiceFactory).WithArgs(axon.Args{"arg1"})
//
// And access Args within a Factory like so.
//    func MyServiceFactory(args axon.Args) axon.Instance {
//      return &myServiceImpl{StringField: args.String(0)} // The arg returned is the string "arg1".
//    }
type Args []interface{}

// Gets a string from the Args. If the passed in index is not a string or does not exist,
// a zero string is returned
func (f Args) String(idx int) (val string) {
	if idx < len(f) {
		if v, ok := f[idx].(string); ok {
			val = v
		}
	}
	return
}

// Gets a int from the Args. If the passed in index is not a int or does not exist,
// a zero int is returned
func (f Args) Int(idx int) (val int) {
	if idx < len(f) {
		if v, ok := f[idx].(int); ok {
			val = v
		}
	}
	return
}

// Gets a int64 from the Args. If the passed in index is not a int64 or does not exist,
// a zero int64 is returned
func (f Args) Int64(idx int) (val int64) {
	if idx < len(f) {
		if v, ok := f[idx].(int64); ok {
			val = v
		}
	}
	return
}

// Gets a uint from the Args. If the passed in index is not a uint or does not exist,
// a zero uint is returned
func (f Args) UInt(idx int) (val uint) {
	if idx < len(f) {
		if v, ok := f[idx].(uint); ok {
			val = v
		}
	}
	return
}

// Gets a float32 from the Args. If the passed in index is not a float32 or does not exist,
// a zero float32 is returned
func (f Args) Float32(idx int) (val float32) {
	if idx < len(f) {
		if v, ok := f[idx].(float32); ok {
			val = v
		}
	}
	return
}

// Gets a float64 from the Args. If the passed in index is not a float64 or does not exist,
// a zero float64 is returned
func (f Args) Float64(idx int) (val float64) {
	if idx < len(f) {
		if v, ok := f[idx].(float64); ok {
			val = v
		}
	}
	return
}

// Gets a bool from the Args. If the passed in index is not a bool or does not exist,
// false is returned
func (f Args) Bool(idx int) (val bool) {
	if idx < len(f) {
		if v, ok := f[idx].(bool); ok {
			val = v
		}
	}
	return
}

// Gets an interface{} from the Args. If the passed in index does not exist,
// a zero interface{} is returned
func (f Args) Any(idx int) (val interface{}) {
	if idx < len(f) {
		val = f[idx]
	}
	return
}

// A collection of Packages. The Binder is the top-level definition of what the Injector will store and how it will store
// it. The Bindings provided by this Binder's Packages will not be evaluated until an injector.Get*(string) method is
// called.
type Binder interface {
	// The slice of Packages that this Binder stores.
	Packages() []Package
}

type binderImpl struct {
	packages []Package
}

func (b *binderImpl) Packages() []Package {
	return b.packages
}

// Create a new Binder from a series of Packages. If no Packages are passed in, the Injector will essentially have no
// Bindings defined and will be functionally useless but no error will be returned.
//  axon.NewBinder(axon.NewPackage(...), axon.NewPackage(...))
func NewBinder(packages ...Package) Binder {
	return &binderImpl{packages: packages}
}

type atomicProvider struct {
	provider Provider
	once     sync.Once
}
