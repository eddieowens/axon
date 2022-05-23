package depgraph

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DoubleMapTestSuite struct {
	suite.Suite
}

func (d *DoubleMapTestSuite) TestAddAndGet() {
	// -- Given
	//
	given := NewDoubleMap[int]()

	// -- When
	//
	given.Add("123", 123)

	// -- Then
	//
	d.Equal(123, given.Get("123"))
}

func (d *DoubleMapTestSuite) TestCircularDependency() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)
	given.Add("3", 3)
	given.AddDependencies("1", "2", "3")
	given.AddDependencies("2", "3")

	// -- When
	//
	given.AddDependencies("3", "1")

	// -- Then
	//
	d.ElementsMatch([]string{"2", "3"}, given.GetDependencies("1"))
	d.ElementsMatch([]string{"3"}, given.GetDependencies("2"))
	d.ElementsMatch([]string{"1"}, given.GetDependencies("3"))

	d.ElementsMatch([]string{"3"}, given.GetDependents("1"))
	d.ElementsMatch([]string{"1"}, given.GetDependents("2"))
	d.ElementsMatch([]string{"2", "1"}, given.GetDependents("3"))
}

func (d *DoubleMapTestSuite) TestRemoveDependencies() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)
	given.AddDependencies("1", "2")

	// -- When
	//
	given.RemoveDependencies("1")

	// -- Then
	//
	d.Empty(given.GetDependencies("1"))
}

func (d *DoubleMapTestSuite) TestAddDependenciesNotExists() {
	// -- Given
	//
	given := NewDoubleMap[int]()

	// -- When
	//
	given.AddDependencies("1", "2")

	// -- Then
	//
	d.Empty(given.GetDependencies("1"))
}

func (d *DoubleMapTestSuite) TestAddDependenciesDepNotExists() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)

	// -- When
	//
	given.AddDependencies("1", "2")

	// -- Then
	//
	d.Empty(given.GetDependencies("1"))
}

func (d *DoubleMapTestSuite) TestRange() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	count := 0

	// -- When
	//
	given.Range(func(key any, val int) bool {
		count++
		return true
	})

	// -- Then
	//
	d.Equal(1, count)
}

func (d *DoubleMapTestSuite) TestRangeDependencies() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)
	given.Add("3", 3)
	given.AddDependencies("1", "2")
	count := 0

	// -- When
	//
	given.RangeDependencies("1", func(_ any, _ int, _ DepMap[any, int]) bool {
		count++
		return true
	})

	// -- Then
	//
	d.Equal(1, count)
}

func (d *DoubleMapTestSuite) TestRangeDependents() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)
	given.Add("3", 3)
	given.AddDependencies("1", "2")
	count := 0

	// -- When
	//
	given.RangeDependents("2", func(_ any, _ int, _ DepMap[any, int]) bool {
		count++
		return true
	})

	// -- Then
	//
	d.Equal(1, count)
}

func (d *DoubleMapTestSuite) TestRemove() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)
	given.AddDependencies("1", "2")

	// -- When
	//
	given.Remove("1")

	// -- Then
	//
	d.Zero(given.Get("1"))
	d.Empty(given.GetDependencies("1"))
	d.Empty(given.GetDependents("1"))
}

func (d *DoubleMapTestSuite) TestFind() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)
	given.Add("2", 2)

	// -- When
	//
	actual := given.Find(func(key any, val int) bool {
		return val == 1
	})

	// -- Then
	//
	d.Equal(1, actual)
}

func (d *DoubleMapTestSuite) TestLookup() {
	// -- Given
	//
	given := NewDoubleMap[int]()
	given.Add("1", 1)

	// -- When
	//
	actual, ok := given.Lookup("1")

	// -- Then
	//
	d.Equal(1, actual)
	d.True(ok)
}

func TestDepGraphTestSuite(t *testing.T) {
	suite.Run(t, new(DoubleMapTestSuite))
}
