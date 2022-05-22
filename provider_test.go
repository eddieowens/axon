package axon

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProviderTestSuite struct {
	suite.Suite
}

func (p *ProviderTestSuite) TestSetValueNil() {
	// -- Given
	//
	given := NewProvider(1)

	// -- When
	//
	err := given.SetValue(nil)

	// -- Then
	//
	p.NoError(err)
}

func (p *ProviderTestSuite) TestSetValueWrongType() {
	// -- Given
	//
	given := NewProvider(1)

	// -- When
	//
	err := given.SetValue("1")

	// -- Then
	//
	p.NoError(err)
}

func TestProviderTestSuite(t *testing.T) {
	suite.Run(t, new(ProviderTestSuite))
}
