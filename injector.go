package axon

import (
	"errors"
	"fmt"
	"github.com/eddieowens/axon/internal/depgraph"
	"github.com/eddieowens/axon/internal/mirror"
	"github.com/eddieowens/axon/opts"
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
	Inject(d any, opts ...opts.Opt[InjectorInjectOpts]) error

	// Add adds the val indexed by a Key. The underlying value for a Key should be a comparable value since the underlying
	// implementation utilizes a map. All calls to Add will overwrite existing values and no checks are done.
	Add(key Key, val any, ops ...opts.Opt[InjectorAddOpts])

	// Get gets a value given a Key. If Get is unable to find the Key, ErrNotFound is returned. The first call to Get will
	// cause the underlying value to be constructed if it is a Factory.
	Get(k Key, o ...opts.Opt[InjectorGetOpts]) (any, error)
}

type InjectorGetOpts struct {
}

type InjectorInjectOpts struct {
	SkipFieldErr bool
}

type InjectorAddOpts struct {
}

// WithSkipFieldErrs allows for the Injector.Inject method to skip over field errors that are encountered when attempting
// to inject values onto a struct. The Injector.Inject method may still return an error, but they will not be due to problems
// encountered on individual fields.
func WithSkipFieldErrs() opts.Opt[InjectorInjectOpts] {
	return func(opts *InjectorInjectOpts) {
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

func (i *injector) Inject(d any, opts ...opts.Opt[InjectorInjectOpts]) error {
	val := reflect.ValueOf(d)
	if val.Kind() != reflect.Ptr || !val.IsValid() {
		return ErrPtrToStruct
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrPtrToStruct
	}

	return i.injectStructWithOpts(Key{}, val, opts...)
}

func (i *injector) Add(key Key, val any, _ ...opts.Opt[InjectorAddOpts]) {
	conVal := newContainerProvider(val)

	if v := key.resolve(i.DepGraph); v != nil {
		i.DepGraph.RemoveDependencies(key)
		v.Invalidate()
	}

	conVal.SetConstructor(func(constructed container[any]) error {
		val := mirror.StripPtrs(constructed.GetReflectValue())

		var err error
		if val.Kind() == reflect.Struct {
			err = i.injectStructWithOpts(key, val)
		}
		return err
	})

	i.DepGraph.Add(key, conVal)
}

func (i *injector) Get(k Key, _ ...opts.Opt[InjectorGetOpts]) (any, error) {
	v := k.resolve(i.DepGraph)
	if v == nil {
		return nil, ErrNotFound
	}

	con, err := v.ProvideContainer()
	if err != nil {
		return nil, err
	}

	return con.GetValue(), nil
}

func (i *injector) injectStructWithOpts(key Key, v reflect.Value, opts ...opts.Opt[InjectorInjectOpts]) error {
	o := &InjectorInjectOpts{}
	for _, v := range opts {
		v(o)
	}

	return i.injectStruct(key, v, *o)
}

func (i *injector) injectStruct(key Key, v reflect.Value, o InjectorInjectOpts) error {
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
	dep := key.resolve(i.DepGraph)
	if dep == nil {
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

	// only applicable to check if MutableValue is used as a field on a struct
	conIsMutable := containerVal.Type().Implements(mutableValueType)
	if conIsMutable {
		err := safeCast[MutableValue](rawVal).SetValue(rawVal)
		if err != nil {
			return fmt.Errorf("failed to set field %s: %w", key.String(), err)
		}
	} else if field.Type().Implements(mutableValueType) {
		if field.IsNil() {
			err := mirror.Instantiate(field)
			if err != nil {
				return fmt.Errorf("%w: field %s is type %s but got type %s", ErrInvalidType, key.String(), field.Type().String(), containerVal.Type().String())
			}
		}

		return getMutableValue(field).SetValue(rawVal)
	}

	err := mirror.Set(field, containerVal)
	if err != nil {
		return fmt.Errorf("%w: field %s is type %s but got type %s", ErrInvalidType, key.String(), field.Type().String(), containerVal.Type().String())
	}
	return nil
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
