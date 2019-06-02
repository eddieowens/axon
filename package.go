package axon

// Builds Bindings which utilize a Factory.
type FactoryBindingBuilder interface {
	// Binds a Key to a provided Factory.
	Factory(factory Factory) FactoryArgsBindingBuilder
}

// Builds Bindings which utilize constants.
type ConstantBindingBuilder interface {
	// Binds the provided Key to a string value.
	String(s string) Binding

	// Binds the provided Key to an int value.
	Int(i int) Binding

	// Binds the provided Key to an int32 value.
	Int32(i int32) Binding

	// Binds the provided Key to an int64 value.
	Int64(i int64) Binding

	// Binds the provided Key to a float32 value.
	Float32(f float32) Binding

	// Binds the provided Key to a float64 value.
	Float64(f float64) Binding

	// Binds the provided Key to a bool value.
	Bool(b bool) Binding

	// Binds the provided Key to a bool value.
	StructPtr(s interface{}) Binding
}

// Builder for a binding to an Injectable Entity (an Instance or a Provider).
type InjectableEntityBindingBuilder interface {
	FactoryBindingBuilder
	ConstantBindingBuilder

	// Bind the provided Key to an Instance.
	Instance(instance Instance) Binding
}

// An intermediary builder for a Binding. Offers no functional impact on the Binding. Is only comprised of prepositions
// to allow the flow of the builder to be more intuitive.
type BindingBuilder interface {
	// Passthrough for defining an Injectable Entity (Provider or Instance).
	To() InjectableEntityBindingBuilder
}

type FactoryArgsBindingBuilder interface {
	// Pass in arguments into your Factory.
	//  Bind("service").To().Factory(ServiceFactory).WithArgs(axon.Args{"my arg"})
	// All arguments passed in here will take precedence over injected values. These will not be
	// overwritten.
	WithArgs(args Args) Binding

	// Specifies that the Factory within the Binding does not have Args.
	WithoutArgs() Binding
}

// A "binder" between a Key and an Injectable Entity (Instance or Provider) within an Injector. A Binding will never hold
// both an Instance and a Provider.
type Binding interface {
	// Retrieves the Instance that this Binding manages. This field is mutually exclusive with the Provider stored
	// in the Binding.
	GetInstance() Instance

	// Retrieves the Provider that this Binding manages. This field is mutually exclusive with the Instance stored
	// in the Binding. The Key tied to the Provider will ultimately be tied to the Instance that the Provider produces,
	// not the Provider itself.
	GetProvider() Provider

	// The key that is bound to the Provider or Instance.
	GetKey() string
}

// Defines a Binding in a Package from a specified key to an Instance or a Provider.
//  axon.NewBinder(axon.NewPackage(
//      axon.Bind("UserService").To().Instance(axon.StructPtr(new(UserService))),
//      axon.Bind("MyString").To().String("hello world"),
//      axon.Bind("DB").To().Factory(DBFactory).WithoutArgs(),
//  ))
func Bind(key string) BindingBuilder {
	m := bindingBuilderImpl{
		key: key,
	}
	return &m
}

// Creates a new Package to pass into your Binder. You can also define a Package as a struct
// by implementing the Package interface like so
//   type ConfigPackage struct {
//   }
//
//   type ConfigPackage struct {
//   }
//
//   func (*ConfigPackage) Entries() []Binding {
//  	 return []Binding{
//       Bind("Dev").To().Factory(TestServiceFactory).WithoutArgs(),
//     }
//   }
//   ...
//   axon.NewBinder(new(ConfigPackage), ...)
func NewPackage(packageBindings ...Binding) Package {
	return &packageImpl{entries: packageBindings}
}

// A group of Bindings that are conceptually tied together. A Package only serves as a way to separate Binding
// definitions into different conceptual chunks. A Package does not provide any functional difference within the Injector
type Package interface {
	// Returns all of the Bindings that this Package stores.
	Bindings() []Binding
}

type packageImpl struct {
	entries []Binding
}

func (m *packageImpl) Bindings() []Binding {
	return m.entries
}

// Defines how to "provide" an Instance to the Injector via a Factory. A Provider's Factory will only ever be called
// once by the Injector. After the Provider's Factory creates your Instance, the injector will stop referring to the
// Provider and will only refer to the Instance that the Provider's Factory produced.
type Provider interface {
	// The Factory that the Provider will use to construct your Instance when called upon by the injector.Get(key string)
	// method.
	GetFactory() Factory

	// Arguments that will be passed into your Factory when it's called.
	GetArgs() Args
}

// Creates a new Provider from the passed in Args and Factory. This function should mainly be used for adding Mock
// Providers into your tests. Rather than calling this function directly within source code, use Bindings instead
// like so
//   axon.Bind("MyService").To().Factory(ServiceFactory).WithoutArgs()
func NewProvider(factory Factory, args ...interface{}) Provider {
	return &providerImpl{
		args:    args,
		factory: factory,
	}
}

type bindingBuilderImpl struct {
	key string
}

func (m *bindingBuilderImpl) To() InjectableEntityBindingBuilder {
	return &injectableEntityBindingBuilderImpl{key: m.key}
}

type packageEntryBuilderImpl struct {
	key     string
	factory Factory
}

func (m *packageEntryBuilderImpl) WithArgs(args Args) Binding {
	return &packageEntryImpl{
		Key: m.key,
		Provider: &providerImpl{
			args:    args,
			factory: m.factory,
		},
	}
}

func (m *packageEntryBuilderImpl) WithoutArgs() Binding {
	return &packageEntryImpl{
		Key: m.key,
		Provider: &providerImpl{
			factory: m.factory,
		},
	}
}

type injectableEntityBindingBuilderImpl struct {
	key string
}

func (m *injectableEntityBindingBuilderImpl) StructPtr(s interface{}) Binding {
	return &packageEntryImpl{Instance: StructPtr(s), Key: m.key}
}

func (m *injectableEntityBindingBuilderImpl) Instance(instance Instance) Binding {
	return &packageEntryImpl{Instance: instance, Key: m.key}
}

func (m *injectableEntityBindingBuilderImpl) Factory(factory Factory) FactoryArgsBindingBuilder {
	return &packageEntryBuilderImpl{key: m.key, factory: factory}
}

func (m *injectableEntityBindingBuilderImpl) String(s string) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: String(s),
	}
}

func (m *injectableEntityBindingBuilderImpl) Int(i int) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Int(i),
	}
}

func (m *injectableEntityBindingBuilderImpl) Int32(i int32) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Int32(i),
	}
}

func (m *injectableEntityBindingBuilderImpl) Int64(i int64) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Int64(i),
	}
}

func (m *injectableEntityBindingBuilderImpl) Float32(f float32) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Float32(f),
	}
}

func (m *injectableEntityBindingBuilderImpl) Float64(f float64) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Float64(f),
	}
}

func (m *injectableEntityBindingBuilderImpl) Bool(b bool) Binding {
	return &packageEntryImpl{
		Key:      m.key,
		Instance: Bool(b),
	}
}

type packageEntryImpl struct {
	Provider Provider
	Instance Instance
	Key      string
}

func (m *packageEntryImpl) GetKey() string {
	return m.Key
}

func (m *packageEntryImpl) GetProvider() Provider {
	return m.Provider
}

func (m *packageEntryImpl) GetInstance() Instance {
	return m.Instance
}

type providerImpl struct {
	factory Factory
	args    Args
}

func (p *providerImpl) GetFactory() Factory {
	return p.factory
}

func (p *providerImpl) GetArgs() Args {
	return p.args
}
