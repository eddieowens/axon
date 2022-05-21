package axon

import (
	"fmt"
	"github.com/eddieowens/axon/internal/depgraph"
	"reflect"
)

type Key struct {
	val       any
	isTypeKey bool
}

func resolve[V any](k Key, d depgraph.DepMap[Key, containerProvider[any]]) containerProvider[any] {
	if k.isTypeKey {
		return d.Find(func(key Key, val containerProvider[any]) bool {
			var ok bool
			if key.isTypeKey {
				_, ok = val.GetComparableValue().(V)
			}
			return ok
		})
	}

	return d.Get(k)
}

func (k Key) String() string {
	return fmt.Sprintf("%v", k.val)
}

func (k Key) IsEmpty() bool {
	return k.val == nil
}

type KeyConstraint interface {
	string
}

func NewKey[V KeyConstraint](val V) Key {
	return Key{val: val}
}

func newReflectKey(v reflect.Value) Key {
	return Key{isTypeKey: true, val: v.Type().String()}
}

func NewTypeKey[V any]() Key {
	return Key{isTypeKey: true, val: reflect.ValueOf(new(V)).Type().Elem().String()}
}
