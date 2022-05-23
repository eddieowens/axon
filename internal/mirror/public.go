package mirror

import (
	"errors"
	"reflect"
)

var (
	ErrIncompatibleTypes = errors.New("incompatible types")
)

// Instantiate if the field is zero, the field is set. If the field can't be set, an error is returned.
func Instantiate(field reflect.Value) error {
	if field.IsZero() {
		val := reflect.New(StripTypePtrs(field.Type()))
		if field.Kind() != reflect.Ptr && field.Kind() != reflect.Interface {
			val = val.Elem()
		}
		return Set(field, val)
	}

	return nil
}

// Set calls dst.Set(src) but won't panic and instead returns an error.
func Set(dst, src reflect.Value) error {
	for !CanSet(dst, src) {
		return ErrIncompatibleTypes
	}

	dst.Set(src)
	return nil
}

// CanSet returns true if the reflect.Value.Set can be called without panic.
func CanSet(dst, src reflect.Value) bool {
	return dst.CanSet() && (dst.Kind() == src.Kind() || (dst.Kind() == reflect.Interface && src.Type().Implements(dst.Type())))
}

// StripTypePtrs returns the core type underneath an arbitrary number of pointers.
func StripTypePtrs(typ reflect.Type) reflect.Type {
	typ, _ = stripTypePtrsCount(typ)
	return typ
}

// StripPtrs returns the core value underneath an arbitrary number of pointers.
func StripPtrs(val reflect.Value) reflect.Value {
	val, _ = stripPtrsCount(val)
	return val
}

func stripTypePtrsCount(typ reflect.Type) (reflect.Type, int) {
	i := 0
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		i++
	}

	return typ, i
}

func stripPtrsCount(val reflect.Value) (reflect.Value, int) {
	i := 0
	for val.Kind() == reflect.Ptr && !val.IsZero() {
		val = val.Elem()
		i++
	}

	return val, i
}
