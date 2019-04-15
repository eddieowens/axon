package tests

import (
	"fmt"
	"github.com/eddieowens/axon"
)

const DepTwoInstanceName = "DepTwo"

type DepTwo interface {
	CallDepTwo() string
}

type DepTwoImpl struct {
	Bool bool `inject:"constantBool"`
}

func (d *DepTwoImpl) CallDepTwo() string {
	return "dep two! " + fmt.Sprint(d.Bool)
}

func DepTwoFactory(args axon.Args) axon.Instance {
	return axon.StructPtr(&DepTwoImpl{Bool: args.Bool(0)})
}
