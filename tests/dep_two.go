package tests

import "axon"

const DepTwoInstanceName = "DepTwo"

type DepTwo interface {
	axon.Instance
	CallDepTwo() string
}

type DepTwoImpl struct {
}

func (*DepTwoImpl) CallDepTwo() string {
	return "dep two!"
}

func (*DepTwoImpl) GetInstanceName() string {
	return DepTwoInstanceName
}
