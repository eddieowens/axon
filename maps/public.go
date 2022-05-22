package maps

type Map[K any, V any] interface {
	Remove(key K)
	Add(key K, val V)
	Get(key K) V
	Lookup(key K) (V, bool)
	Range(r RangeFunc[K, V])
	Find(r FindFunc[K, V]) V
}

type RangeFunc[K any, V any] func(key K, val V) bool

type FindFunc[K any, V any] func(key K, val V) bool
