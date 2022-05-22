package axon

import (
	"errors"
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
	i.Equal(2, actual.Int.Get())
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
	inj.Add(NewTypeKey[MutableValue](&mutable{Int: 2}))
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

func (i *InjectorTestSuite) TestInjectFailedSetValue() {
	// -- Given
	//
	type test struct {
		Mutable MutableValue `inject:",type"`
	}

	inj := NewInjector()
	inj.Add(NewTypeKey[MutableValue](mutableError{}))
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	i.EqualError(err, "failed to set field axon.MutableValue: not found")
}

func (i *InjectorTestSuite) TestInjectTypePtrImpl() {
	// -- Given
	//
	type test struct {
		T testInterface `inject:",type"`
	}

	inj := NewInjector()
	inj.Add(NewTypeKey[testInterface](&testInterfacePtr{Int: 1}))

	actual := new(test)
	expected := &test{T: &testInterfacePtr{Int: 1}}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestInjectTypeValImpl() {
	// -- Given
	//
	type test struct {
		T testInterface `inject:",type"`
	}

	inj := NewInjector()
	inj.Add(NewTypeKey[testInterface](testInterfaceVal{Int: 1}))

	actual := new(test)
	expected := &test{T: testInterfaceVal{Int: 1}}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestInjectTypeAndKey() {
	// -- Given
	//
	type test struct {
		T testInterface `inject:",type"`
		A testInterface `inject:"a"`
	}

	inj := NewInjector()
	inj.Add(NewTypeKey[testInterface](testInterfaceVal{Int: 1}))
	inj.Add(NewKey("a"), &testInterfacePtr{Int: 2})

	actual := new(test)
	expected := &test{T: testInterfaceVal{Int: 1}, A: &testInterfacePtr{Int: 2}}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestInjectMissingField() {
	// -- Given
	//
	type test struct {
		I int `inject:"i"`
	}

	inj := NewInjector()
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	i.EqualError(err, "failed to inject i: not found")
}

func (i *InjectorTestSuite) TestInjectWrongType() {
	// -- Given
	//
	type test struct {
		I int `inject:"i"`
	}

	inj := NewInjector()
	inj.Add(NewKey("i"), "1")

	// -- When
	//
	err := inj.Inject(new(test))

	// -- Then
	//
	i.EqualError(err, "invalid type: field i is type int but got type string")
}

func (i *InjectorTestSuite) TestInjectSkipErr() {
	// -- Given
	//
	type test struct {
		I int    `inject:"i"`
		S string `inject:",type"`
	}

	inj := NewInjector()
	inj.Add(NewKey("i"), "1")
	inj.Add(NewTypeKey[string]("str"))
	expected := &test{
		S: "str",
	}
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual, WithSkipFieldErrs())

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestNonStructMutableValue() {
	// -- Given
	//
	type test struct {
		M mutableMap `inject:"m"`
	}

	inj := NewInjector()

	inj.Add(NewKey("m"), mutableMap{"1": "2"})

	expected := &test{
		M: mutableMap{"1": "21"},
	}
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestStructProvider() {
	// -- Given
	//
	type providedStruct struct {
		I int        `inject:"i"`
		M mutableMap `inject:"m"`
		S int
	}

	type test struct {
		P *Provider[*providedStruct] `inject:"p"`
	}

	inj := NewInjector()

	inj.Add(NewKey("p"), &providedStruct{S: 3})
	inj.Add(NewKey("i"), 2)
	inj.Add(NewKey("m"), mutableMap{"1": "2"})

	expected := &test{P: Provide(&providedStruct{S: 3, I: 2, M: mutableMap{"1": "21"}})}
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestInjectNonPtr() {
	// -- Given
	//
	inj := NewInjector()

	// -- When
	//
	err := inj.Inject(testDep{})

	// -- Then
	//
	i.EqualError(err, "value must be a ptr to a struct")
}

func (i *InjectorTestSuite) TestInjectNonStruct() {
	// -- Given
	//
	inj := NewInjector()
	p := 1

	// -- When
	//
	err := inj.Inject(&p)

	// -- Then
	//
	i.EqualError(err, "value must be a ptr to a struct")
}

func (i *InjectorTestSuite) TestInjectNil() {
	// -- Given
	//
	inj := NewInjector()
	var given *test

	// -- When
	//
	err := inj.Inject(given)

	// -- Then
	//
	i.EqualError(err, "value must be a ptr to a struct")
}

func (i *InjectorTestSuite) TestFailedFactory() {
	// -- Given
	//
	inj := NewInjector()
	inj.Add(NewKey("1"), NewFactory[string](func(inj Injector) (string, error) {
		return "", errors.New("error")
	}))

	// -- When
	//
	actual, err := inj.Get(NewKey("1"))

	// -- Then
	//
	i.EqualError(err, "error")
	i.Nil(actual)
}

func (i *InjectorTestSuite) TestInjectWithFactory() {
	// -- Given
	//
	type test struct {
		S string `inject:"s"`
	}

	inj := NewInjector()
	inj.Add(NewKey("s"), NewFactory[string](func(inj Injector) (string, error) {
		return "1", nil
	}))
	actual := new(test)
	expected := &test{S: "1"}

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	if i.NoError(err) {
		i.Equal(expected, actual)
	}
}

func (i *InjectorTestSuite) TestInjectFailedFactory() {
	// -- Given
	//
	type test struct {
		S string `inject:"s"`
	}

	inj := NewInjector()
	inj.Add(NewKey("s"), NewFactory[string](func(inj Injector) (string, error) {
		return "", errors.New("error")
	}))
	actual := new(test)

	// -- When
	//
	err := inj.Inject(actual)

	// -- Then
	//
	i.EqualError(err, "failed to get field s: error")
}

func (i *InjectorTestSuite) TestUnsettableField() {
	// -- Given
	//
	type test struct {
		unset int `inject:"s"`
	}

	inj := NewInjector()
	inj.Add(NewKey("s"), 1)

	// -- When
	//
	err := inj.Inject(new(test))

	// -- Then
	//
	i.EqualError(err, "invalid field: field s is not settable")
}

func (i *InjectorTestSuite) TestFailedMutableValue() {
	// -- Given
	//
	type test struct {
		M MutableValue `inject:"m"`
	}

	inj := NewInjector()
	inj.Add(NewKey("m"), mutable{})

	// -- When
	//
	err := inj.Inject(new(test))

	// -- Then
	//
	i.EqualError(err, "invalid type: field m is type axon.MutableValue but got type axon.mutable")
}

type testInterface interface {
	testSigil()
}

type testInterfaceVal struct {
	Int int
}

func (t testInterfaceVal) testSigil() {}

type testInterfacePtr struct {
	Int int
}

func (t *testInterfacePtr) testSigil() {}

type mutable struct {
	Int int
}

type mutableError struct {
}

func (m mutableError) SetValue(_ any) error {
	return ErrNotFound
}

type mutableMap map[string]string

func (m mutableMap) SetValue(val any) error {
	in, ok := val.(mutableMap)
	if !ok {
		return ErrInvalidType
	}

	for k, v := range in {
		m[k] = v + "1"
	}
	return nil
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
