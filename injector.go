package axon

import (
	"fmt"
	"reflect"
	"sync"
)

// An injector is, in essence, a map of a string key to either an [Instance](#instance) or a
// Provider. When a user calls GetInstance(string) they are providing
// the key into the map which will either simply return the instance, or construct it via the
// Provider.Factory field.
//
// This is the main object you will interact with to get whatever you place into the injector via a Binder
type Injector interface {
	// Gets a specific Instance within the Injector. If the Instance does not yet exist (created via a Provider) it is constructed here.
	// If the key is not found a panic occurs.
	GetInstance(key string) Instance

	// Add a Provider to the Injector so the next time the Instance is being retrieved, the following
	// Provider will be used to build it.
	AddProvider(key string, provider *Provider)

	// Add an Instance to the Injector. If the Instance with the referent (created by Instance.GetInstanceName()) already
	// exists, it is overwritten with the provided Instance.
	//
	// Every time this is called, all dependencies of whatever was added will be rebuilt on subsequent GetInstance(string)
	// calls.
	AddInstance(instance Instance)

	// Same as AddInstance but allows for the manual override of the Instance's referent within
	// the injector.
	AddInstanceWithKey(key string, instance Instance)
}

// An Instance is an interface with a single method; GetInstanceName() string. This
// method is what is used to identify the Instance when it is within the Injector.
// For example
//   type TestService interface {
//       axon.Instance // All implementations of this service are now an axon.Instance
//       DoServiceThings() string
//   }
//
//   type TestServiceImpl struct {
//   }
//
//   func (*TestServiceImpl) GetInstanceName() string {
//       return "testService"
//   }
//
//   func (*TestServiceImpl) DoServiceThings() string {
//       return "I'm a service!"
//   }
// After adding the TestService to your Binder via a BinderEntry, all calls to
// injector.GetInstance("testService") will return TestServiceImpl. To override the implementation define
// another TestService called TestServiceImplTwo
//   type TestServiceImplTwo struct {
//   }
//
//   func (*TestServiceImplTwo) GetInstanceName() string {
//       return "testService" // The exact same name as TestServiceImpl.GetInstanceName
//   }
//
//   func (*TestServiceImplTwo) DoServiceThings() string {
//       return "I'm the second service!"
//   }
// Now call
//   injector.AddInstance(new(TestServiceImplTwo))
// to override the key "testService" in the injector
// with the TestServiceImplTwo instance. Every subsequent call to injector.GetInstance("testService) will
// now return TestServiceImplTwo.
//
// This is also true for all calls to Instances that depend on a TestService. If you were to define another
// service that depends on TestService like so
//   type HigherLevelService interface {
//       axon.Instance
//   }
//
//   type HigherLevelServiceImpl struct {
//       TestService TestService inject:"testService" // TestService dependency
//   }
//
//   func (TestServiceImpl) GetInstanceName() string {
//       return "higherLevelService"
//   }
//
//   func HigherLevelServiceFactory() axon.Instance {
//       return new(HigherLevelServiceImpl)
//   }
// And add HigherLevelService to your Binder through a Provider using the
// HigherLevelServiceFactory function, whenever injector.GetInstance("higherLevelService") is called,
// it will populate the TestService field with the TestServiceImplTwo implementation of TestService.
//
// Any field without the inject:"INSTANCE_NAME" tag will not be populated by the injector
type Instance interface {
	// The referent used to retrieve this Instance from the Injector
	GetInstanceName() string
}

type injectorImpl struct {
	// Instances that were provided by the Binder. These have already been
	// constructed by GetInstance, or provided by the user in the Binder
	binderInstances binderInstances

	// A map of InstanceName to a Provider. Used for constructing the Instance when
	// GetInstance is called
	atomicProviderMap map[string]*atomicProvider

	// A lock for all writes to the state stored in the injector
	lock sync.Mutex

	// A map of InstanceNames to InstanceNames that depend on them
	dependencyMap dependencyMap
}

type dependencyMap map[string][]string

func (i *injectorImpl) AddInstanceWithKey(key string, instance Instance) {
	if instance == nil {
		return
	}
	i.lock.Lock()
	defer i.lock.Unlock()
	i.binderInstances[key] = instance
	i.clearInstanceDeps(key)
}

func (i *injectorImpl) AddInstance(instance Instance) {
	if instance == nil {
		return
	}
	i.lock.Lock()
	defer i.lock.Unlock()
	i.binderInstances[instance.GetInstanceName()] = instance
	i.clearInstanceDeps(instance.GetInstanceName())

}

func (i *injectorImpl) AddProvider(key string, provider *Provider) {
	if provider == nil {
		return
	}
	i.lock.Lock()
	defer i.lock.Unlock()
	i.atomicProviderMap[key] = &atomicProvider{
		provider: provider,
		once:     sync.Once{},
	}
	i.clearInstanceDeps(key)
}

func (i *injectorImpl) clearInstanceDeps(instanceName string) {
	for _, v := range i.dependencyMap[instanceName] {
		delete(i.binderInstances, v)
		i.atomicProviderMap[v].once = sync.Once{}
	}
}

func (i *injectorImpl) GetInstance(key string) Instance {
	instance := i.binderInstances[key]
	if instance == nil {
		ap := i.atomicProviderMap[key]
		if ap == nil {
			panic(fmt.Sprintf("unknown instance %s", key))
		}
		ap.once.Do(func() {
			if ap.provider.Args == nil {
				instance = ap.provider.Factory(nil)
			} else {
				instance = ap.provider.Factory(ap.provider.Args)
			}
			v := reflect.ValueOf(instance).Elem()
			for j := 0; j < v.NumField(); j++ {
				depInstanceName := v.Type().Field(j).Tag.Get("inject")
				if depInstanceName != "" {
					depInstance := i.GetInstance(depInstanceName)
					if v.Field(j).CanSet() {
						v.Field(j).Set(reflect.ValueOf(depInstance))
						i.dependencyMap[depInstanceName] = append(i.dependencyMap[depInstanceName], ap.provider.InstanceName)
					}
				}
			}
			i.lock.Lock()
			defer i.lock.Unlock()
			i.binderInstances[key] = instance
		})
	}
	return instance
}

func newInjector(binder Binder) Injector {
	bi, ap, dm := hydrateInjector(binder)
	return &injectorImpl{
		binderInstances:   bi,
		atomicProviderMap: ap,
		lock:              sync.Mutex{},
		dependencyMap:     dm,
	}
}

func hydrateInjector(binder Binder) (binderInstances, map[string]*atomicProvider, dependencyMap) {
	bi := make(binderInstances)
	ap := make(map[string]*atomicProvider)
	dm := make(dependencyMap)
	for _, v := range binder {
		if v.Instance != nil {
			bi[v.Instance.GetInstanceName()] = v.Instance
		} else {
			if v.Provider == nil {
				panic("all binder entries must have either a Val or a InstanceFactory")
			}
			if v.Provider.InstanceName == "" {
				panic("if only a provider factory is used, an InstanceName must be provided in the binder entry")
			}
			ap[v.Provider.InstanceName] = &atomicProvider{
				once:     sync.Once{},
				provider: v.Provider,
			}
		}
	}
	return bi, ap, dm
}
