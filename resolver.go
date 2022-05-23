package axon

import (
	"github.com/eddieowens/axon/maps"
)

type StorageGetter interface {
	Find(f maps.FindFunc[any, containerProvider[any]]) containerProvider[any]
	Get(k any) containerProvider[any]
}

type Resolver interface {
	resolve(storage StorageGetter) containerProvider[any]
}
