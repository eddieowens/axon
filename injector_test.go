package axon

import (
	"github.com/eddieowens/axon/internal/depgraph"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InjectorTestSuite struct {
	suite.Suite
}

func (i *InjectorTestSuite) TestAdd() {
	// -- Given
	//
	inj := &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}

	// -- When
	//
	inj.Add(NewKey("key"), 23)

	// -- Then
	//
	actual, _ := inj.Get(NewKey("key"))
	i.Equal(23, actual)
}

func (i *InjectorTestSuite) TestAddStruct() {
	// -- Given
	//
	inj := &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}
	p := 2
	inj.Add(NewKey("I"), 1)
	inj.Add(NewKey("PtrI"), &p)
	inj.Add(NewKey("S"), "s")
	inj.Add(NewKey("Test"), new(test))
	inj.Add(NewKey("dep"), new(testDep))

	expected := &test{
		I:    1,
		PtrI: &p,
		Dep: &testDep{
			S: "s",
		},
	}

	// -- When
	//
	actual, err := inj.Get(NewKey("Test"))

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
		i.ElementsMatch([]Key{NewKey("I"), NewKey("PtrI"), NewKey("dep")}, inj.DepGraph.GetDependencies(NewKey("Test")))
		i.ElementsMatch([]Key{NewKey("S")}, inj.DepGraph.GetDependencies(NewKey("dep")))
		i.ElementsMatch([]Key{}, inj.DepGraph.GetDependents(NewKey("Test")))
		i.ElementsMatch([]Key{NewKey("Test")}, inj.DepGraph.GetDependents(NewKey("dep")))
	}
}

func (i *InjectorTestSuite) TestInject() {
	// -- Given
	//
	inj := &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}
	p := 2
	inj.Add(NewKey("I"), 1)
	inj.Add(NewKey("PtrI"), &p)
	inj.Add(NewKey("S"), "s")
	inj.Add(NewKey("dep"), new(testDep))

	expected := &test{
		I:    1,
		PtrI: &p,
		Dep: &testDep{
			S: "s",
		},
	}
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
		i.ElementsMatch([]Key{NewKey("S")}, inj.DepGraph.GetDependencies(NewKey("dep")))
		i.ElementsMatch([]Key{}, inj.DepGraph.GetDependents(NewKey("dep")))
		i.ElementsMatch([]Key{NewKey("dep")}, inj.DepGraph.GetDependents(NewKey("S")))
	}
}

func (i *InjectorTestSuite) TestAddNewDependency() {
	// -- Given
	//
	inj := &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}
	p := 2
	inj.Add(NewKey("I"), 1)
	inj.Add(NewKey("PtrI"), &p)
	inj.Add(NewKey("S"), "s")
	inj.Add(NewKey("dep"), new(testDep))
	inj.Add(NewKey("Test"), new(test))
	_, _ = inj.Get(NewKey("Test"))

	// -- When
	//
	inj.Add(NewKey("dep"), new(testDep))

	// -- Then
	//
	i.ElementsMatch([]Key{NewKey("I"), NewKey("PtrI"), NewKey("dep")}, inj.DepGraph.GetDependencies(NewKey("Test")))
	i.ElementsMatch([]Key{NewKey("Test")}, inj.DepGraph.GetDependents(NewKey("dep")))
	i.ElementsMatch([]Key{}, inj.DepGraph.GetDependencies(NewKey("dep")))
	i.ElementsMatch([]Key{}, inj.DepGraph.GetDependents(NewKey("S")))
}

func (i *InjectorTestSuite) TestSettingMutableValue() {
	// -- Given
	//
	type test struct {
		Int   *Provider[int]      `inject:"im"`
		Strct *Provider[*testDep] `inject:"sm"`
	}

	inj := &injector{
		DepGraph: depgraph.NewDoubleMap[containerProvider[any]](),
	}

	inj.Add(NewKey("S"), "str")
	inj.Add(NewKey("sm"), &testDep{})
	inj.Add(NewKey("im"), 1)

	actual := new(test)
	err := inj.Inject(actual)
	i.NoError(err)

	expected := &test{
		Int: Provide(1),
		Strct: Provide(&testDep{
			S: "str",
		}),
	}

	i.Equal(expected, actual)

	// -- When
	//
	actual.Int.Set(2)

	// -- Then
	//
	expected.Int = Provide(2)
	i.Equal(expected, actual)
}

func (i *InjectorTestSuite) TestMutableImplementation() {
	// -- Given
	//
	type test struct {
		Mutable *mutable `inject:"mutable"`
	}

	inj := NewInjector()
	inj.Add(NewKey("mutable"), &mutable{Int: 2})

	actual := new(test)
	expected := &test{Mutable: &mutable{Int: 1}}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestMutableImplementationWrongType() {
	// -- Given
	//
	type test struct {
		Mutable *mutable `inject:"mutable"`
	}

	inj := NewInjector()
	inj.Add(NewKey("mutable"), mutable{Int: 2})
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	i.EqualError(err, ErrInvalidType.Error())
}

func (i *InjectorTestSuite) TestInjectType() {
	// -- Given
	//
	type test struct {
		Mutable MutableValue `inject:",type"`
	}

	inj := NewInjector()
	inj.Add(NewTypeKey[MutableValue](), &mutable{Int: 2})
	actual := new(test)
	expected := &test{Mutable: &mutable{Int: 1}}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

type mutable struct {
	Int int
}

func (m *mutable) SetValue(val any) error {
	in, ok := val.(*mutable)
	if !ok {
		return ErrInvalidType
	}
	m.Int = in.Int - 1
	return nil
}

type testDep struct {
	S string `inject:"S"`
}

type test struct {
	I    int      `inject:"I"`
	PtrI *int     `inject:"PtrI"`
	Dep  *testDep `inject:"dep"`
}

func TestInjectorTestSuite(t *testing.T) {
	suite.Run(t, new(InjectorTestSuite))
}
