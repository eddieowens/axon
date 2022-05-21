package axon

type Factory interface {
	Build(inj Injector) (any, error)
}

type internalFactory interface {
	GetZeroValue() any
}

func NewFactory[T any](f FactoryFunc[T]) Factory {
	return &factory[T]{
		Val:         *new(T),
		FactoryFunc: f,
	}
}

type factory[T any] struct {
	FactoryFunc FactoryFunc[T]

	// A zero value of the type the factory builds.
	Val T
}

func (f *factory[T]) GetZeroValue() any {
	return f.Val
}

func (f *factory[T]) Build(inj Injector) (any, error) {
	return f.FactoryFunc.Build(inj)
}

type FactoryFunc[T any] func(inj Injector) (T, error)

func (f FactoryFunc[T]) Build(inj Injector) (any, error) {
	return f(inj)
}
