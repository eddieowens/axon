package axon

import (
	"errors"
	"fmt"
	"github.com/eddieowens/axon/internal/depgraph"
	"reflect"
	"strings"
)

var (
	InjectTag          = "inject"
	InjectTagValueType = "type"
)

var (
	ErrPtrToStruct  = errors.New("value must be a ptr to a struct")
	ErrNotFound     = errors.New("not found")
	ErrInvalidType  = errors.New("invalid type")
	ErrInvalidField = errors.New("invalid field")
)

// Injector allows for the storage, retrieval, and construction of objects.
type Injector interface {
	// Inject injects all fields on a struct that are tagged with the InjectTag from the Injector. d must be a pointer to
	// a struct and the fields that are tagged must be public. If the InjectTag is not present on the struct or if the
	// value is already set, it will not be injected.
	//
	// errors returned:
	//
	// ErrNotFound: If a field has an InjectTag but is not found in the Injector.
	//
	// ErrInvalidField: If the field is not settable.
	//
	// ErrInvalidType: If the type in the Injector is not the same as the type that is being injected.
	//
	// ErrPtrToStruct: If d is not a pointer to a struct.
	//
	// All errors should be checked with errors.Is as they may be wrapped.
	Inject(d any, opts ...InjectOpt) error

	// Add adds the val indexed by a Key. The underlying value for a Key should be a comparable value since the underlying
	// implementation utilizes a map. All calls to Add will overwrite existing values and no checks are done.
	Add(key Key, val any)

	// Get gets a value given a Key. If Get is unable to find the Key, ErrNotFound is returned. The first call to Get will
	// cause the underlying value to be constructed if it is a Factory.
	Get(k Key) (any, error)

	getGraph() depgraph.DoubleMap[any, containerProvider[any]]
}

type InjectOpt func(opts *InjectOpts)

type InjectOpts struct {
	SkipFieldErr bool
}

// WithSkipFieldErrs allows for the Injector.Inject method to skip over field errors that are encountered when attempting
// to inject values onto a struct. The Injector.Inject method may still return an error, but they will not be due to problems
// encountered on individual fields.
func WithSkipFieldErrs() InjectOpt {
	return func(opts *InjectOpts) {
		opts.SkipFieldErr = true
	}
}

// NewInjector constructs a new Injector.
func NewInjector() Injector {
	return &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}
}

// MutableValue allows for any implementation to control how a value is being set within the Injector.
//
// Normally the Injector sets an injected value via reflect.ValueOf(instance).Set(value) but if you want to have some
// instance managed by the Injector be mutated in a specific way, you can implement MutableValue with a pointer receiver.
type MutableValue interface {

	// SetValue is called when the Injector is looking to inject a value to the implementor. val may not be the same expected
	// type or could also be nil. All implementations should do type checks as well as nil checks.
	//
	// Because this causes a mutation of the implementor, this method should always be implemented via a pointer receiver.
	SetValue(val any) error
}

var mutableValueType = reflect.TypeOf((*MutableValue)(nil)).Elem()

type injector struct {
	DepGraph depgraph.DoubleMap[any, containerProvider[any]]
}

func (i *injector) getGraph() depgraph.DoubleMap[any, containerProvider[any]] {
	return i.DepGraph
}

func (i *injector) Inject(d any, opts ...InjectOpt) error {
	val := reflect.ValueOf(d)
	if val.Kind() != reflect.Ptr || !val.IsValid() {
		return ErrPtrToStruct
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrPtrToStruct
	}

	return i.injectStructWithOpts(Key{}, d, val, opts...)
}

func (i *injector) Add(key Key, val any) {
	if val == nil {
		return
	}

	conVal := newContainerProvider(val)

	if v, ok := i.DepGraph.Lookup(key); ok {
		i.DepGraph.RemoveDependencies(key)
		v.Invalidate()
	}

	conVal.SetConstructor(func(constructed container[any]) error {
		val := stripPtrs(constructed.GetReflectValue())

		var err error
		if val.Kind() == reflect.Struct {
			err = i.injectStructWithOpts(key, constructed.GetValue(), val)
		}
		return err
	})

	i.DepGraph.Add(key, conVal)
}

func (i *injector) Get(k Key) (any, error) {
	v := i.DepGraph.Get(k)
	if v == nil {
		return nil, ErrNotFound
	}

	con, err := v.ProvideContainer()
	if err != nil {
		return nil, err
	}

	return con.GetValue(), nil
}

func (i *injector) injectStructWithOpts(key Key, rawVal any, v reflect.Value, opts ...InjectOpt) error {
	o := &InjectOpts{}
	for _, v := range opts {
		v(o)
	}

	mut := safeCast[MutableValue](rawVal)
	if mut != nil {
		v = stripPtrs(reflect.ValueOf(mut))
		// if the underlying value of the MutableValue is not a struct, we need to grab it, and set it.
		if v.Kind() != reflect.Struct {
			if !key.IsEmpty() {
				con, err := i.resolveValue(key)
				if err != nil {
					return err
				}

				conVal := con.GetValue()
				return mut.SetValue(conVal)
			}
			return nil
		}
	}

	return i.injectStruct(key, v, *o)
}

func (i *injector) injectStruct(key Key, v reflect.Value, o InjectOpts) error {
	for j := 0; j < v.NumField(); j++ {
		err := i.injectStructField(key, v.Field(j), v.Type().Field(j))
		if err != nil {
			if o.SkipFieldErr {
				continue
			}
			return err
		}
	}

	return nil
}

func (i *injector) injectStructField(key Key, field reflect.Value, strctField reflect.StructField) error {
	depInjectTag := strctField.Tag.Get(InjectTag)
	depKey := resolveKey(depInjectTag, field)
	if !depKey.IsEmpty() {
		con, err := i.resolveValue(depKey)
		if err != nil {
			return err
		}

		err = setReflectVal(field, con, depKey)
		if err != nil {
			return err
		}

		i.DepGraph.AddDependencies(key, depKey)
	}
	return nil
}

func (i *injector) resolveValue(key Key) (container[any], error) {
	dep, ok := i.DepGraph.Lookup(key)
	if !ok {
		return nil, fmt.Errorf("failed to inject %s: %w", key.String(), ErrNotFound)
	}

	con, err := dep.ProvideContainer()
	if err != nil {
		return nil, fmt.Errorf("failed to get field %s: %w", key.String(), err)
	}

	return con, nil
}

func getMutableValue(field reflect.Value) (out MutableValue) {
	if field.CanInterface() {
		i := field.Interface()
		if i != nil {
			out, _ = i.(MutableValue)
		}
	}
	return
}

func safeCast[T any](val any) (out T) {
	if val != nil {
		out, _ = val.(T)
	}
	return
}

func setReflectVal(field reflect.Value, container container[any], key Key) error {
	if !field.CanSet() {
		return fmt.Errorf("%w: field %s is not settable", ErrInvalidField, key.String())
	}

	rawVal := container.GetValue()
	containerVal := container.GetReflectValue()

	if containerVal.Type().Implements(mutableValueType) {
		return safeCast[MutableValue](rawVal).SetValue(rawVal)
	}

	if field.Kind() != containerVal.Kind() {
		return fmt.Errorf("%w: field %s is type %s while the value in the injector is type %s", ErrInvalidType, key.String(), field.Type().String(), containerVal.Type().String())
	}

	field.Set(containerVal)
	return nil
}

func stripPtrs(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr && !val.IsZero() {
		val = val.Elem()
	}

	return val
}

func stripTypePtrs(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func resolveKey(tag string, field reflect.Value) Key {
	parsed := parseTag(tag)
	if field.IsZero() && parsed != nil {
		if parsed.Name != "" {
			return NewKey(parsed.Name)
		} else if parsed.InjectType {
			return newReflectKey(field)
		}
	}

	return Key{}
}

func parseTag(tag string) *injectTag {
	if tag == "" {
		return nil
	}
	out := &injectTag{}

	tagSplit := strings.Split(tag, ",")
	if len(tagSplit) > 1 {
		for _, v := range tagSplit[1:] {
			if strings.TrimSpace(v) == InjectTagValueType {
				out.InjectType = true
			}
		}
	}

	out.Name = strings.TrimSpace(tagSplit[0])

	return out
}

type injectTag struct {
	// The first field in the InjectTag. The name will always take the highest precedence when resolving the Key for the
	// Injector.
	Name string

	// Corresponds to the InjectTagValueType field of the tag. If that value is present in the InjectTag, the type of the
	// dependency is injected rather than a specific key.
	InjectType bool
}
