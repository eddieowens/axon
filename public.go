package axon

import "fmt"

var DefaultInjector = NewInjector()

type GetOpt func(opts *GetOpts)

type InjectableKey interface {
	Key | string
}

func WithKey[V InjectableKey](k V) GetOpt {
	return func(opts *GetOpts) {
		opts.Key = injectableKeyToKey(k)
	}
}

type GetOpts struct {
	Key Key
}

func MustGet[V any](opts ...GetOpt) V {
	out, err := Get[V](opts...)
	if err != nil {
		panic(err)
	}
	return out
}

func Get[V any](opts ...GetOpt) (V, error) {
	return InjectorGet[V](DefaultInjector, opts...)
}

func InjectorGet[V any](inj Injector, opts ...GetOpt) (out V, err error) {
	o := &GetOpts{}
	for _, v := range opts {
		v(o)
	}

	if !o.Key.IsEmpty() {
		val, err := inj.Get(o.Key)
		if err != nil {
			return out, err
		}

		if out, ok := val.(V); ok {
			return out, nil
		} else {
			return out, fmt.Errorf("%w: expected %s key to be type %T but got %T", ErrInvalidType, o.Key.String(), out, val)
		}
	} else {
		val := inj.getGraph().Find(func(key any, val containerProvider[any]) bool {
			_, ok := val.GetComparableValue().(V)
			return ok
		})
		if val == nil {
			return out, ErrNotFound
		}

		con, err := val.ProvideContainer()
		if err != nil {
			return out, err
		}

		conVal := con.GetValue()
		if out, ok := conVal.(V); ok {
			return out, nil
		} else {
			return out, fmt.Errorf("expected type %T but got %T: %w", out, val, ErrInvalidType)
		}
	}
}

type AddOpt func(opts *AddOpts)

type AddOpts struct {
	AsType bool
}

func Add[K InjectableKey](key K, val any, opts ...AddOpt) {
	InjectAdd(DefaultInjector, key, val, opts...)
}

func InjectAdd[K InjectableKey](inj Injector, key K, val any, opts ...AddOpt) {
	o := &AddOpts{}
	for _, v := range opts {
		v(o)
	}

	inj.Add(injectableKeyToKey(key), newContainerProvider(val))
}

func Provide[T any](val T) *Provider[T] {
	return &Provider[T]{val: val}
}

func injectableKeyToKey[V InjectableKey](key V) Key {
	var k Key
	switch typ := any(key).(type) {
	case string:
		k = NewKey(typ)
	case Key:
		k = typ
	}
	return k
}
