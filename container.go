package axon

import (
	"reflect"
)

type container[T any] interface {
	GetValue() T
	GetReflectValue() reflect.Value
	SetValue(val T)
}

func newContainer[T any](v T) container[T] {
	return &containerImpl[T]{
		Value: v,
	}
}

type containerImpl[T any] struct {
	Value        T
	ReflectValue *reflect.Value
}

func (c *containerImpl[T]) SetValue(val T) {
	c.Value = val
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
