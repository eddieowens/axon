package axon

import (
	"sync"
)

// Factory for creating an Instance. The args field are provided via the Provider listed
// within a BinderEntry
type InstanceFactory func(args Args) Instance

// Args that are passed into an InstanceFactory upon creation. Alias with methods attached
// for easier getting of args.
type Args []interface{}

// Gets a string from the Args. If the passed in index is not a string or does not exist,
// an zero string is returned
func (f Args) String(idx int) (val string) {
	if idx < len(f) {
		if v, ok := f[idx].(string); ok {
			val = v
		}
	}
	return
}

// Gets a int from the Args. If the passed in index is not a int or does not exist,
// an zero int is returned
func (f Args) Int(idx int) (val int) {
	if idx < len(f) {
		if v, ok := f[idx].(int); ok {
			val = v
		}
	}
	return
}

// Gets a int64 from the Args. If the passed in index is not a int64 or does not exist,
// an zero int64 is returned
func (f Args) Int64(idx int) (val int64) {
	if idx < len(f) {
		if v, ok := f[idx].(int64); ok {
			val = v
		}
	}
	return
}

// Gets a uint from the Args. If the passed in index is not a uint or does not exist,
// an zero uint is returned
func (f Args) UInt(idx int) (val uint) {
	if idx < len(f) {
		if v, ok := f[idx].(uint); ok {
			val = v
		}
	}
	return
}

// Gets a float32 from the Args. If the passed in index is not a float32 or does not exist,
// an zero float32 is returned
func (f Args) Float32(idx int) (val float32) {
	if idx < len(f) {
		if v, ok := f[idx].(float32); ok {
			val = v
		}
	}
	return
}

// Gets a float64 from the Args. If the passed in index is not a float64 or does not exist,
// an zero float64 is returned
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
// an zero interface{} is returned
func (f Args) Any(idx int) (val interface{}) {
	if idx < len(f) {
		val = f[idx]
	}
	return
}

// A Provider is used within a BinderEntry as an alternative to an Instance. The
// Provider's role is to _provide_ a factory to the Injector for creating your Instance
// when it is requested via injector.GetInstance(string).
//
// Should my BinderEntry take a Provider or an Instance
//
// Generally, if your Instance requires another Instance, a factory function should be defined, and
// added to the BinderEntry via a Provider. If you construct your Instance outside of the BinderEntry,
// passing in all of the required dependencies, the Injector will not know what is required by your Instance
// and will not be able to pass mocks of dependencies in for example.
//
// This is why Instance BinderEntry's should be used for leaf dependencies and Providers should be used for all
// other dependencies.
type Provider struct {
	// The function used to create your Instance when it is called for by injector.GetInstance(InstanceName string)
	Factory InstanceFactory

	// Arguments that will be passed into your Factory when it's created
	Args Args

	// The InstanceName must be provided here because your Instance will not be created
	// until asked from the injector via injector.GetInstanceName(NAME). Until your Instance
	// is constructed via the Factory, the injector does not know its referent so it must be
	// provided here
	InstanceName string
}

// Defines how an Injector will provide an Instance when it is requested via
// injector.GetInstance(string). It is ultimately a slice of BinderEntry.
//
// The definition of a Binder should be scoped to a function so that it may be easily used within tests like so
//   func CreateInjector() Injector {
//       binder := []BinderEntry{
//           {
//               Provider: &Provider{
//                   Factory:      TestServiceFactory,
//                   InstanceName: "testService",
//               },
//           },
//           {
//               Instance: &TestServiceDependencyImpl{},
//           },
//       }
//
//       return NewInjector(binder)
//   }
// Now this CreateInjector() function can be called within your tests to create the injector.
type Binder []BinderEntry

type binderInstances map[string]Instance

// Defines how the Injector will either build, or return your Instance when it is requested via
// injector.GetInstance(string). A BinderEntry can be given a user-created Instance OR
// a Provider but not both. If both are added to the BinderEntry, the Instance will be
// used.
//
// To note
//
// If an Instance is passed into the BinderEntry rather than a Provider, all inject
// tags on the Instance's fields will be ignored as they are being provided by a user. This defeats the
// purpose of this library so only provide an Instance to the BinderEntry if it is a leaf in the dependency
// graph.
type BinderEntry struct {
	// A struct which defines a factory to create your Instance and any arguments that should be bound
	// to the InstanceFactory function. The InstanceName field in the BinderEntry is required if this field is
	// to be used
	Provider *Provider

	// The raw Instance used in the Injector. The output of the Instance.GetInstanceName() function will
	// be used as the Instance's referent.
	//
	// All Instances that are passed into the Binder using this field will not be able to have their dependencies
	// overwritten as the Instance was constructed outside of the Injector. If you wish to let the Injector be
	// able to overwrite dependencies via AddInstance(), use the Provider instead.
	Instance Instance
}

type atomicProvider struct {
	provider *Provider
	once     sync.Once
}
