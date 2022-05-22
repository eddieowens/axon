package axon

import (
	"fmt"
	"reflect"
)

// Key the key type for the Injector.
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

// NewTypeKey returns a Key using the type V as well as the passed in value val.
func NewTypeKey[V any](val V) (Key, V) {
	return newTypeKey[V](), val
}

// NewTypeKeyFactory rather than tying a particular type to a value like NewTypeKey, this func ties a type to a Factory.
func NewTypeKeyFactory[V any](val Factory) (Key, Factory) {
	return newTypeKey[V](), val
}

// NewTypeKeyProvider rather than tying a particular type to a value like NewTypeKey, this func ties a type to a provider.
func NewTypeKeyProvider[V any](val V) (Key, *Provider[V]) {
	return newTypeKey[V](), NewProvider(val)
}

// newTypeKey returns a Key based on the type of V.
func newTypeKey[V any]() Key {
	return Key{
		isTypeKey: true,
		val:       reflect.ValueOf(new(V)).Type().Elem().String(),
	}
}

func newReflectKey(v reflect.Value) Key {
	return Key{isTypeKey: true, val: v.Type().String()}
}
