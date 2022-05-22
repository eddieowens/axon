package axon

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PublicTestSuite struct {
	suite.Suite
}

func (p *PublicTestSuite) BeforeTest(_, _ string) {
	DefaultInjector = NewInjector()
}

func (p *PublicTestSuite) TestGet() {
	// -- Given
	//
	Add("", "2")
	Add("1", 3)

	// -- When
	//
	actual, err := Get[string]()

	// -- Then
	//
	if p.NoError(err) {
		p.Equal(3, MustGet[int](WithKey("1")))
		p.Equal("2", actual)
	}
}

func (p *PublicTestSuite) TestGetWithKey() {
	// -- Given
	//
	Add("", "2")
	Add("1", 3)

	// -- When
	//
	actual, err := Get[int](WithKey(NewKey("1")))

	// -- Then
	//
	if p.NoError(err) {
		p.Equal(3, actual)
	}
}

func (p *PublicTestSuite) TestGetWrongType() {
	// -- Given
	//
	Add("1", "2")

	// -- When
	//
	actual, err := Get[int](WithKey("1"))

	// -- Then
	//
	if p.EqualError(err, "invalid type: expected 1 key to be type int but got string") {
		p.Equal(0, actual)
	}
}

func (p *PublicTestSuite) TestGetMissing() {
	// -- When
	//
	actual, err := Get[int]()

	// -- Then
	//
	p.EqualError(err, "not found")
	p.Zero(actual)
}

func (p *PublicTestSuite) TestWithFactory() {
	// -- Given
	//
	given := NewFactory(func(inj Injector) (int, error) {
		return 1, nil
	})

	Add("2", given)

	// -- When
	//
	actual, err := Get[int](WithKey("2"))

	// -- Then
	//
	if p.NoError(err) {
		p.Equal(1, actual)
	}
}

func (p *PublicTestSuite) TestWithFactoryWrongType() {
	// -- Given
	//
	given := NewFactory(func(inj Injector) (string, error) {
		return "123", nil
	})

	Add("2", given)

	// -- When
	//
	actual, err := Get[int](WithKey("2"))

	// -- Then
	//
	p.EqualError(err, "invalid type: expected 2 key to be type int but got string")
	p.Zero(actual)
}

func (p *PublicTestSuite) TestWithFactoryError() {
	// -- Given
	//
	given := NewFactory(func(inj Injector) (string, error) {
		return "", errors.New("error")
	})

	Add("2", given)

	// -- When
	//
	actual, err := Get[string](WithKey("2"))

	// -- Then
	//
	p.EqualError(err, "error")
	p.Zero(actual)
}

func (p *PublicTestSuite) TestMustGetKeyPanics() {
	p.Panics(func() {
		MustGet[int]()
	})
}

func TestPublicTestSuite(t *testing.T) {
	suite.Run(t, new(PublicTestSuite))
}
