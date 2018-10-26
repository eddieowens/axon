package axon

import (
	"fmt"
	"sync"
)

// Allows for getting and updating singletons
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
	// Every time this is called, the dependency graph will be cleared and rebuilt
	AddInstance(instance Instance)

	// Same as AddInstance but allows for the manual override of the Instance's referent within
	// the injector.
	AddInstanceWithKey(key string, instance Instance)
}

// A single entry within the Injector e.g.
//   type TestService interface {
//     axon.Instance
//     DoTestStuff() string
//   }
type Instance interface {
	// The referent used to retrieve this Instance from the Injector
	GetInstanceName() string
}

type injectorImpl struct {
	binderInstances   binderInstances
	binder            Binder
	atomicProviderMap map[string]*atomicProvider
	lock              sync.Mutex
}

func (i *injectorImpl) AddInstanceWithKey(key string, instance Instance) {
	i.lock.Lock()
	defer i.lock.Unlock()
	bi, ap := hydrateInjector(i.binder)
	i.binderInstances = bi
	i.atomicProviderMap = ap
	i.binderInstances[key] = instance
}

func (i *injectorImpl) AddInstance(instance Instance) {
	i.lock.Lock()
	defer i.lock.Unlock()
	bi, ap := hydrateInjector(i.binder)
	i.binderInstances = bi
	i.atomicProviderMap = ap
	i.binderInstances[instance.GetInstanceName()] = instance
}

func (i *injectorImpl) AddProvider(key string, provider *Provider) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.atomicProviderMap[key] = &atomicProvider{
		provider: provider,
		once:     sync.Once{},
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
			instance = ap.provider.Factory(i, ap.provider.Args...)
			i.lock.Lock()
			defer i.lock.Unlock()
			i.binderInstances[key] = instance
		})
	}
	return instance
}

func newInjector(binder Binder) Injector {
	bi, ap := hydrateInjector(binder)
	return &injectorImpl{
		binderInstances:   bi,
		atomicProviderMap: ap,
		lock:              sync.Mutex{},
	}
}

func hydrateInjector(binder Binder) (binderInstances, map[string]*atomicProvider) {
	bi := make(binderInstances)
	ap := make(map[string]*atomicProvider)
	for _, v := range binder {
		if v.Instance != nil {
			if v.Name == "" {
				bi[v.Instance.GetInstanceName()] = v.Instance
			} else {
				bi[v.Name] = v.Instance
			}
		} else {
			if v.Provider == nil {
				panic("all binder entries must have either a Val or a InstanceFactory")
			}
			if v.Name == "" {
				panic("if only a provider factory is used, a name must be provided in the binder entry")
			}
			ap[v.Name] = &atomicProvider{
				once:     sync.Once{},
				provider: v.Provider,
			}
		}
	}
	return bi, ap
}
