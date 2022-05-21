package depgraph

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type DepGraphTestSuite struct {
	suite.Suite
}

func (d *DepGraphTestSuite) TestAddAndGet() {
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

func (d *DepGraphTestSuite) TestCircularDependency() {
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

func TestDepGraphTestSuite(t *testing.T) {
	suite.Run(t, new(DepGraphTestSuite))
}
