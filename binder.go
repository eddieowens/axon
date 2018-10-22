package axon

import (
	"sync"
)

// Factory for creating an Instance. The args field are provided via the Provider listed
// within a BinderEntry
type InstanceFactory func(injector Injector, args ...interface{}) Instance

// Defines the InstanceFactory and Args that will be used to build the Instance
// when it is being retrieved from the Injector
type Provider struct {
	Factory InstanceFactory
	Args    []interface{}
}

// A simple list of BinderEntry
type Binder []BinderEntry

type binderInstances map[string]Instance

// Defines how the Injector will create, and reference an Instance.
type BinderEntry struct {
	// Overrides the referent for the Instance which is listed in the Instance.GetInstanceName() method.
	// Required if no Instance is provided to the BinderEntry
	Name string

	// A struct which defines a factory to create your Instance and any arguments that should be bound
	// to the InstanceFactory function. The Name field in the BinderEntry is required if this field is
	// to be used
	Provider *Provider

	// The raw Instance used in the Injector. The output of the Instance.GetInstanceName() function will
	// be used as the Instance's referent
	Instance Instance
}

type atomicProvider struct {
	provider *Provider
	once     sync.Once
}
