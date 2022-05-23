package axon

import (
	"github.com/eddieowens/axon/opts"
	"sync"
)

// Provider allows values to be mutated in real-time in a thread-safe manner. Providers should be used when you have a
// source of data which you want to be updated in multiple places at once even after leaving the Injector. Values returned
// by Provider.Get should never be stored as that defeats the purpose of the Provider altogether. Instead, Provider.Get should
// be called every time one wants to read the data which the Provider provides.
//
// For example
//
//    Add("one", 1)
//    one := MustGet[int]("one")
//    fmt.Println(one) // prints 1
//    Add("one", 2)
//    fmt.Println(one) // still prints 1
//
// To allow for dynamic value updates, you can use a provider.
//
//    Add("one", NewProvider(1))
//    one := MustGet[*Provider[int]]("one")
//    fmt.Println(one.Get()) // prints 1
//    Add("one", NewProvider(2))
//    fmt.Println(one.Get()) // prints 2
type Provider[T any] struct {
	val  T
	lock sync.RWMutex
}

func (p *Provider[T]) Set(val T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.val = val
}

// SetValue for a Provider, this supports a val of either T or *Provider[T].
func (p *Provider[T]) SetValue(val any) error {
	if val == nil {
		return nil
	}

	v, ok := val.(T)
	if !ok {
		prov, ok := val.(*Provider[T])
		if ok {
			v, ok = any(prov.val).(T)
			if !ok {
				return ErrInvalidType
			}
		}
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

	// GetValue returns the underlying value for the containerProvider. This value may be the wrapped value or a zero value
	// if not yet constructed.
	GetValue() T

	// IsInstantiated returns true if ProvideContainer has ever been called, false otherwise.
	IsInstantiated() bool
}

type OnConstructFunc[T any] func(constructed container[T]) error

func newContainerProvider(inj Injector, val any) containerProvider[any] {
	p := &containerProviderImpl[any]{
		Value:    val,
		Injector: inj,
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
	Value        T
	Container    container[T]
	Factory      Factory
	Once         sync.Once
	Instantiated bool
	Injector     Injector
	OnConstruct  OnConstructFunc[T]
}

func (p *containerProviderImpl[T]) GetValue() T {
	return p.Value
}

func (p *containerProviderImpl[T]) IsInstantiated() bool {
	return p.Instantiated
}

func (p *containerProviderImpl[T]) SetConstructor(constructor OnConstructFunc[T]) {
	p.OnConstruct = constructor
}

func (p *containerProviderImpl[T]) ProvideContainer() (container[T], error) {
	var err error
	p.Once.Do(func() {
		val := p.Value
		if p.Container == nil {
			kt := newKeyTracker(p.Injector)
			if p.Factory != nil {
				var v any
				v, err = p.Factory.Build(kt)
				if err != nil {
					return
				}
				val = v.(T)
			}
			p.Container = newContainer(val, kt.keysGotten...)

			if p.OnConstruct != nil {
				err = p.OnConstruct(p.Container)
			}
		}
	})
	if err != nil {
		return nil, err
	}

	p.Instantiated = true
	return p.Container, nil
}

func (p *containerProviderImpl[T]) Invalidate() {
	p.Once = sync.Once{}
}

func newKeyTracker(i Injector) *keyTracker {
	return &keyTracker{
		Injector:   i,
		keysGotten: make([]Key, 0),
	}
}

type keyTracker struct {
	Injector
	keysGotten []Key
}

func (t *keyTracker) Get(k Key, _ ...opts.Opt[InjectorGetOpts]) (any, error) {
	t.keysGotten = append(t.keysGotten, k)
	return t.Injector.Get(k)
}
