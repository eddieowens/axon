package axon

import (
	"fmt"
	"reflect"
)

type Key struct {
	val       any
	isTypeKey bool
}

func (k Key) resolve(storage StorageGetter) containerProvider[any] {
	return storage.Get(k)
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

func NewTypeKey[V any](val V) (Key, V) {
	return Key{
		isTypeKey: true,
		val:       reflect.ValueOf(new(V)).Type().Elem().String(),
	}, val
}

func newReflectKey(v reflect.Value) Key {
	return Key{isTypeKey: true, val: v.Type().String()}
}
