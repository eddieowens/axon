package axon

import (
	"fmt"
	"reflect"
	"sync"
)

// An Injector is, in essence, a map of a string Key to an Injectable Entity (Instance or a Provider). When a user calls
// injector.Get(key string) they are providing the key into the map which will either simply return the instance,
// or construct it via the Provider's Factory.
//
// This is the main object you will interact with to get whatever you place into the injector via a Binder.
type Injector interface {
	// Gets a specific Instance within the Injector. If the Instance does not yet exist (created via a Provider) it is
	// constructed here. If the key is not found, nil is returned.
	Get(key string) Instance

	// Get a struct Instance from the injector. This will always return a pointer to whatever struct
	// was passed in via the Binding. If not found, nil is returned.
	GetStructPtr(key string) interface{}

	// Get a bool Instance from the Injector. If not found, false is returned.
	GetBool(key string) bool

	// Get a int Instance from the Injector. If not found, 0 is returned.
	GetInt(key string) int

	// Get a string Instance from the Injector. If not found, "" is returned.
	GetString(key string) string

	// Get a float32 Instance from the Injector. If not found, 0 is returned.
	GetFloat32(key string) float32

	// Get a float64 Instance from the Injector. If not found, 0 is returned.
	GetFloat64(key string) float64

	// Get a int64 Instance from the Injector. If not found, 0 is returned.
	GetInt64(key string) int64

	// Get a int32 Instance from the Injector. If not found, 0 is returned.
	GetInt32(key string) int32

	// Add an Instance to the Injector. If the Instance with the specified key already exists, it is replaced.
	//
	// Every time this is called, all dependencies of the Instance will be rebuilt on subsequent Get*(string)
	// calls.
	//
	// WARNING: Do not use this method at runtime within source code. This method is for TESTING purposes in order to
	// provide mocks without having to define a completely separate Binder. If you use this within source code, it will
	// work within a synchronous environment but has undefined behavior in an asynchronous environment.
	Add(key string, instance Instance)

	// Add a Provider to the Injector. If the Instance with the specified key already exists, it is replaced.
	//
	// Every time this is called, all dependencies of the Instance will be rebuilt on subsequent Get*(string)
	// calls.
	//
	// WARNING: Do not use this method at runtime within source code. This method is for TESTING purposes in order to
	// provide mocks without having to define a completely separate Binder. If you use this within source code, it will
	// work within a synchronous environment but has undefined behavior in an asynchronous environment.
	AddProvider(key string, provider Provider)
}

type injectorImpl struct {
	// Instances that were provided by the Binder. These have already been
	// constructed by GetValue, or provided by the user in the Binder
	binderMap injectorMap

	// A map of Key to a Provider. Used for constructing the Instance when
	// GetValue is called
	atomicProviderMap map[string]*atomicProvider

	// A map of InstanceNames to InstanceNames that depend on them
	dependencyMap dependencyMap
}

func (i *injectorImpl) AddProvider(key string, provider Provider) {
	if provider == nil {
		return
	}
	i.atomicProviderMap[key] = newAtomicProvider(provider)
	i.clearInstanceDeps(key)
}

func (i *injectorImpl) GetStructPtr(key string) interface{} {
	inst := i.Get(key)
	if inst != nil {
		return inst.GetStructPtr()
	}
	return inst
}

func (i *injectorImpl) GetBool(key string) bool {
	inst := i.Get(key)
	if inst == nil {
		return false
	}
	return inst.GetBool()
}

func (i *injectorImpl) GetInt(key string) int {
	inst := i.Get(key)
	if inst == nil {
		return 0
	}
	return inst.GetInt()
}

func (i *injectorImpl) GetString(key string) string {
	inst := i.Get(key)
	if inst == nil {
		return ""
	}
	return inst.GetString()
}

func (i *injectorImpl) GetFloat32(key string) float32 {
	inst := i.Get(key)
	if inst == nil {
		return 0
	}
	return inst.GetFloat32()
}

func (i *injectorImpl) GetFloat64(key string) float64 {
	inst := i.Get(key)
	if inst == nil {
		return 0
	}
	return inst.GetFloat64()
}

func (i *injectorImpl) GetInt64(key string) int64 {
	inst := i.Get(key)
	if inst == nil {
		return 0
	}
	return inst.GetInt64()
}

func (i *injectorImpl) GetInt32(key string) int32 {
	inst := i.Get(key)
	if inst == nil {
		return 0
	}
	return inst.GetInt32()
}

type dependencyMap map[string][]string

func (i *injectorImpl) Add(instanceName string, instance Instance) {
	if instance == nil {
		return
	}
	i.binderMap[instanceName] = newManagedInstance(instance, false)
	i.clearInstanceDeps(instanceName)
}

func (i *injectorImpl) Get(key string) Instance {
	mInst := i.binderMap[key]
	instance := mInst.instance
	if instance == nil {
		ap := i.atomicProviderMap[key]
		if ap == nil {
			return nil
		}
		ap.once.Do(func() {
			if ap.provider.GetArgs() == nil {
				instance = ap.provider.GetFactory()(i, nil)
			} else {
				instance = ap.provider.GetFactory()(i, ap.provider.GetArgs())
			}
			if instance.GetKind() == reflect.Struct {
				i.instantiateStructValue(key, instance)
			}
			i.binderMap[key] = newManagedInstance(instance, true)
		})
	} else if !mInst.isInstantiated && instance.GetKind() == reflect.Struct {
		i.instantiateStructValue(key, instance)
		i.binderMap[key] = newManagedInstance(instance, true)
	}
	return instance
}

func (i *injectorImpl) instantiateStructValue(key string, instance Instance) {
	v := instance.getReflectValue().Elem()
	for j := 0; j < v.NumField(); j++ {
		depKey := v.Type().Field(j).Tag.Get("inject")
		if depKey != "" {
			depInstance := i.Get(depKey)
			if depInstance == nil {
				panic(fmt.Sprintf("failed to inject %s into %s as it was not created.", depKey, v.String()))
			}
			if depInstance.GetKind() == reflect.Struct || isZero(v.Field(j)) {
				if v.Field(j).CanSet() {
					v.Field(j).Set(reflect.ValueOf(depInstance.GetValue()))
					i.dependencyMap[depKey] = append(i.dependencyMap[depKey], key)
				}
			}
		}
	}
}

func isZero(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func newAtomicProvider(provider Provider) *atomicProvider {
	return &atomicProvider{
		once:     sync.Once{},
		provider: provider,
	}
}

func (i *injectorImpl) clearInstanceDeps(key string) {
	for _, v := range i.dependencyMap[key] {
		delete(i.binderMap, v)
		i.atomicProviderMap[v].once = sync.Once{}
		i.clearInstanceDeps(v)
	}
}

func newInjector(binder Binder) Injector {
	bi, ap, dm := hydrateInjector(binder)
	return &injectorImpl{
		binderMap:         bi,
		atomicProviderMap: ap,
		dependencyMap:     dm,
	}
}

type managedInstance struct {
	isInstantiated bool
	instance       Instance
}

func newManagedInstance(instance Instance, isInstantiated bool) managedInstance {
	return managedInstance{
		isInstantiated: isInstantiated,
		instance:       instance,
	}
}

type injectorMap map[string]managedInstance

func hydrateInjector(binder Binder) (injectorMap, map[string]*atomicProvider, dependencyMap) {
	bi := make(injectorMap)
	ap := make(map[string]*atomicProvider)
	dm := make(dependencyMap)
	for _, m := range binder.Modules() {
		for _, v := range m.Bindings() {
			if v.GetInstance() != nil {
				bi[v.GetKey()] = newManagedInstance(v.GetInstance(), false)
			} else if v.GetProvider() != nil {
				ap[v.GetKey()] = newAtomicProvider(v.GetProvider())
			}
		}
	}
	return bi, ap, dm
}
