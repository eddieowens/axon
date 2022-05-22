// Package axon is a simple and generic-friendly DI library.
package axon

import (
	"fmt"
	"github.com/eddieowens/axon/opts"
)

// DefaultInjector acts as a global-level Injector for all operations. If you want to create your own Injector, use NewInjector.
var DefaultInjector = NewInjector()

// InjectableKey is a type constraint for the supported keys within the Injector.
type InjectableKey interface {
	Key | string
}

// Inject injects all fields in val marked with the InjectTag using the DefaultInjector. If any errors are encountered
// when trying to inject values to val, an error is returned. val must be a ptr to a struct.
func Inject[V any](val V, opt ...opts.Opt[InjectorInjectOpts]) error {
	return DefaultInjector.Inject(val, opt...)
}

// WithKey specifies an InjectableKey to Get. If not specified, all Get funcs return the specified generic type.
//    Add("mykey", 1)
//    mykey := MustGet[int](WithKey("mykey"))
//    fmt.Println(mykey) // prints 1
func WithKey[V InjectableKey](k V) opts.Opt[GetOpts] {
	return func(opts *GetOpts) {
		opts.Key = injectableKeyToKey(k)
	}
}

type GetOpts struct {
	// A specific Key to get from the Injector.
	Key Key
}

// MustGet same as InjectorGet but panics if an error is encountered.
func MustGet[V any](opts ...opts.Opt[GetOpts]) V {
	out, err := Get[V](opts...)
	if err != nil {
		panic(err)
	}
	return out
}

// Get same as InjectorGet but uses the DefaultInjector.
func Get[V any](opts ...opts.Opt[GetOpts]) (V, error) {
	return InjectorGet[V](DefaultInjector, opts...)
}

// InjectorGet calls Injector.Get but also adds support for generics. This adds some convenience in type casting
//
//    Add(axon.NewTypeKey(1)) // anytime someone asks for the type "int" the value "1" is injected.
//    myInt, err := InjectorGet[int]()
//    fmt.Println(myInt) // prints "1"
func InjectorGet[V any](inj Injector, ops ...opts.Opt[GetOpts]) (out V, err error) {
	o := opts.ApplyOpts(&GetOpts{}, ops...)

	key := o.Key
	if key.IsEmpty() {
		key, _ = NewTypeKey[V](out)
	}

	val, err := inj.Get(key)
	if err != nil {
		return out, err
	}

	if out, ok := val.(V); ok {
		return out, nil
	} else {
		return out, fmt.Errorf("%w: expected %s key to be type %T but got %T", ErrInvalidType, o.Key.String(), out, val)
	}
}

type AddOpts struct {
}

// Add same as InjectAdd but uses the DefaultInjector.
func Add[K InjectableKey](key K, val any, opts ...opts.Opt[AddOpts]) {
	InjectAdd(DefaultInjector, key, val, opts...)
}

// InjectAdd adds a value into the inj using a key. If InjectAdd is called on a pre-existing value, it is overwritten.
func InjectAdd[K InjectableKey](inj Injector, key K, val any, _ ...opts.Opt[AddOpts]) {
	inj.Add(injectableKeyToKey(key), val)
}

func NewProvider[T any](val T) *Provider[T] {
	return &Provider[T]{val: val}
}

func injectableKeyToKey[V InjectableKey](key V) Key {
	var k Key
	if s, ok := any(key).(string); ok {
		if s == "" {
			k, _ = NewTypeKey[V](key)
		} else {
			k = NewKey(s)
		}
	} else {
		k = any(key).(Key)
	}
	return k
}
