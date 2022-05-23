package axon

import (
	"reflect"
)

type container[T any] interface {
	GetValue() T
	GetReflectValue() reflect.Value
	// GetExternalDependencies returns Keys that may be required that aren't explicitly listed on the container's value e.g. dependencies grabbed in a Factory.
	GetExternalDependencies() []Key
}

func newContainer[T any](v T, externalDeps ...Key) container[T] {
	return &containerImpl[T]{
		Value:                v,
		ExternalDependencies: externalDeps,
	}
}

type containerImpl[T any] struct {
	Value                T
	ReflectValue         *reflect.Value
	ExternalDependencies []Key
}

func (c *containerImpl[T]) GetExternalDependencies() []Key {
	return c.ExternalDependencies
}

func (c *containerImpl[T]) GetValue() T {
	return c.Value
}

func (c *containerImpl[T]) GetReflectValue() reflect.Value {
	if c.ReflectValue == nil {
		v := reflect.ValueOf(c.Value)
		c.ReflectValue = &v
		return v
	}
	return *c.ReflectValue
}
