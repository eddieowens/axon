package axon

import (
	"sync"
)

// Provider allows values to be mutated in real-time in a thread-safe manner. Providers should be used when you have a
// source of data which you want to be updated in multiple places at once even after leaving the Injector. Values returned
// by Provider.Get should never be stored as that defeats the purpose of the Provider altogether. Instead, Provider.Get should
// be called every time one wants to read the data which the Provider provides.
type Provider[T any] struct {
	val  T
	lock sync.RWMutex
}

func (p *Provider[T]) Set(val T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.val = val
}

func (p *Provider[T]) SetValue(val any) error {
	if val == nil {
		return nil
	}

	v, ok := val.(T)
	if !ok {
		return nil
	}

	p.Set(v)
	return nil
}

func (p *Provider[T]) Get() T {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.val
}

type containerProvider[T any] interface {
	ProvideContainer() (container[T], error)
	Invalidate()
	SetConstructor(constructor OnConstructFunc[T])
}

type OnConstructFunc[T any] func(constructed container[T]) error

func newContainerProvider(val any) containerProvider[any] {
	p := &containerProviderImpl[any]{
		Value: val,
	}

	fact, ok := val.(Factory)
	if ok {
		p.Factory = fact
	}

	internal, ok := val.(internalFactory)
	if ok {
		p.Value = internal.GetZeroValue()
	}

	return p
}

type containerProviderImpl[T any] struct {
	Value       T
	Container   container[T]
	Factory     Factory
	Once        sync.Once
	Injector    Injector
	OnConstruct OnConstructFunc[T]
}

func (p *containerProviderImpl[T]) SetConstructor(constructor OnConstructFunc[T]) {
	p.OnConstruct = constructor
}

func (p *containerProviderImpl[T]) ProvideContainer() (container[T], error) {
	var err error
	p.Once.Do(func() {
		val := p.Value
		if p.Container == nil {
			if p.Factory != nil {
				var v any
				v, err = p.Factory.Build(p.Injector)
				if err != nil {
					return
				}
				val = v.(T)
			}
			p.Container = newContainer(val)

			if p.OnConstruct != nil {
				err = p.OnConstruct(p.Container)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	return p.Container, nil
}

func (p *containerProviderImpl[T]) Invalidate() {
	p.Once = sync.Once{}
}
