// An Instance of data stored in the Injector. This is what will be retrieved from and added to the Injector.
package axon

import "reflect"

type Instance interface {
	// Get the reflect Kind of the Instance. To note, this reflect.Kind is consistent within
	// the axon codebase but may not be consistent with the raw Go type. For instance, a Struct
	// type is actually a reflect.Ptr.
	GetKind() reflect.Kind

	// The pointer to the struct that the Instance is managing.
	GetStructPtr() interface{}

	// The int that the Instance is managing.
	GetInt() int

	// The int32 that the Instance is managing.
	GetInt32() int32

	// The int64 that the Instance is managing.
	GetInt64() int64

	// The bool that the Instance is managing.
	GetBool() bool

	// The string that the Instance is managing.
	GetString() string

	// The float32 that the Instance is managing.
	GetFloat32() float32

	// The float64 that the Instance is managing.
	GetFloat64() float64

	// The raw value that the Instance is managing.
	GetValue() interface{}

	getReflectValue() reflect.Value
}

// An Instance of any type. This should be used as sparingly as possible as
// more reflection is needed to manage it (which is slower). Use this if none of
// the other Instance types can fit your needs.
func Any(instance interface{}) Instance {
	return &instanceImpl{
		Kind:  reflect.Interface,
		Value: instance,
	}
}

// An Instance of an int type.
func Int(instance int) Instance {
	return &instanceImpl{
		Kind:  reflect.Int,
		Int:   instance,
		Value: instance,
	}
}

// An Instance of a struct ptr type. This type MUST be a ptr to a struct value.
// If it is not, a panic will occur.
// Good:
//   func MyStructFactory(_ axon.Args) instance.Instance {
//     instance.StructPtr(new(MyStruct))
//   }
//
// Bad:
//   func MyStructFactory(_ axon.Args) instance.Instance {
//     instance.StructPtr(MyStruct{})
//   }
//
// Will be stored within the Injector simply as &MyStruct{}
func StructPtr(instance interface{}) Instance {
	v := reflect.ValueOf(instance)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("You are attempting to add a raw struct value. Only pointers to structs are permitted.")
	}
	return &instanceImpl{
		ReflectValue: v,
		Value:        instance,
		Kind:         reflect.Struct,
	}
}

// An Instance of an int32 type.
func Int32(instance int32) Instance {
	return &instanceImpl{
		Int32: instance,
		Kind:  reflect.Int32,
		Value: instance,
	}
}

// An Instance of an int64 type.
func Int64(instance int64) Instance {
	return &instanceImpl{
		Int64: instance,
		Kind:  reflect.Int64,
		Value: instance,
	}
}

// An Instance of a bool type.
func Bool(instance bool) Instance {
	return &instanceImpl{
		Bool:  instance,
		Kind:  reflect.Bool,
		Value: instance,
	}
}

// An Instance of a string type.
func String(instance string) Instance {
	return &instanceImpl{
		String: instance,
		Kind:   reflect.String,
		Value:  instance,
	}
}

// An Instance of a float32 type.
func Float32(instance float32) Instance {
	return &instanceImpl{
		Float32: instance,
		Kind:    reflect.Float32,
		Value:   instance,
	}
}

// An Instance of a float64 type.
func Float64(instance float64) Instance {
	return &instanceImpl{
		Float64: instance,
		Kind:    reflect.Float64,
		Value:   instance,
	}
}

type instanceImpl struct {
	Kind         reflect.Kind
	Value        interface{}
	Int          int
	Int32        int32
	Int64        int64
	Bool         bool
	String       string
	Float32      float32
	Float64      float64
	ReflectValue reflect.Value
}

func (v *instanceImpl) getReflectValue() reflect.Value {
	return v.ReflectValue
}

func (v *instanceImpl) GetValue() interface{} {
	return v.Value
}

func (v *instanceImpl) GetKind() reflect.Kind {
	return v.Kind
}

func (v *instanceImpl) GetStructPtr() interface{} {
	return v.Value
}

func (v *instanceImpl) GetInt() int {
	return v.Int
}

func (v *instanceImpl) GetInt32() int32 {
	return v.Int32
}

func (v *instanceImpl) GetInt64() int64 {
	return v.Int64
}

func (v *instanceImpl) GetBool() bool {
	return v.Bool
}

func (v *instanceImpl) GetString() string {
	return v.String
}

func (v *instanceImpl) GetFloat32() float32 {
	return v.Float32
}

func (v *instanceImpl) GetFloat64() float64 {
	return v.Float64
}
