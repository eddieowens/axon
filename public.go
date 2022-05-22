package axon

import (
	"fmt"
	"github.com/eddieowens/axon/opts"
)

var DefaultInjector = NewInjector()

type InjectableKey interface {
	Key | string
}

func WithKey[V InjectableKey](k V) opts.Opt[GetOpts] {
	return func(opts *GetOpts) {
		opts.Key = injectableKeyToKey(k)
	}
}

type GetOpts struct {
	Key Key
}

func MustGet[V any](opts ...opts.Opt[GetOpts]) V {
	out, err := Get[V](opts...)
	if err != nil {
		panic(err)
	}
	return out
}

func Get[V any](opts ...opts.Opt[GetOpts]) (V, error) {
	return InjectorGet[V](DefaultInjector, opts...)
}

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
	AsType bool
}

func Add[K InjectableKey](key K, val any, opts ...opts.Opt[AddOpts]) {
	InjectAdd(DefaultInjector, key, val, opts...)
}

func InjectAdd[K InjectableKey](inj Injector, key K, val any, _ ...opts.Opt[AddOpts]) {
	inj.Add(injectableKeyToKey(key), val)
}

func Provide[T any](val T) *Provider[T] {
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
