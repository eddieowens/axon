package depgraph

type MutableRangeFunc[K any, V any] func(key K, val V, p DepMap[K, V]) bool
type RangeFunc[K any, V any] func(key K, val V) bool
type FindFunc[K any, V any] func(key K, val V) bool

type DoubleMap[K any, V any] interface {
	DepMap[K, V]

	// GetDependencies recursively searches for all values that the key depends on.
	GetDependencies(key K) []K

	// GetDependents recursively searches for all values that depend on a key.
	GetDependents(key K) []K

	RangeDependencies(key K, r MutableRangeFunc[K, V])

	RangeDependents(key K, r MutableRangeFunc[K, V])

	RemoveDependencies(key K)
}

type Map[K any, V any] interface {
	Remove(key K)
	Add(key K, val V)
	Get(key K) V
	Lookup(key K) (V, bool)
	Range(r RangeFunc[K, V])
	Find(r FindFunc[K, V]) V
}

type DepMap[K any, V any] interface {
	Map[K, V]
	AddDependencies(key K, keys ...K)
}

// NewDoubleMap creates a DoubleMap comprised of two maps rather than a directed graph. This DoubleMap provides constant time lookups
// but more limited searching capabilities.
func NewDoubleMap[V any]() DoubleMap[any, V] {
	return &doubleMap[V]{
		Dependents:   map[any]Set[any]{},
		Dependencies: map[any]Set[any]{},
		Vals:         map[any]V{},
	}
}

type doubleMap[V any] struct {
	// Keys that are depended upon by a set of keys.
	Dependents map[any]Set[any]

	// Keys that depend upon a set of keys.
	Dependencies map[any]Set[any]
	Vals         map[any]V
}

func (m *doubleMap[V]) Find(r FindFunc[any, V]) (v V) {
	for k, v := range m.Vals {
		if r(k, v) {
			return v
		}
	}
	return
}

func (m *doubleMap[V]) RemoveDependencies(key any) {
	for _, dep := range m.Dependents {
		dep.Remove(key)
	}
	delete(m.Dependencies, key)
}

func (m *doubleMap[V]) AddDependencies(key any, keys ...any) {
	_, exists := m.Vals[key]
	if !exists {
		return
	}

	dependencies := getOrInit(m.Dependencies, key)
	for _, valKey := range keys {
		_, exists := m.Vals[valKey]
		if !exists {
			continue
		}

		dependencies.Add(valKey)
		dependents := getOrInit(m.Dependents, valKey)
		dependents.Add(key)
	}
}

func (m *doubleMap[V]) Range(r RangeFunc[any, V]) {
	for k, v := range m.Vals {
		ok := r(k, v)
		if !ok {
			return
		}
	}
}

func (m *doubleMap[V]) RangeDependencies(key any, r MutableRangeFunc[any, V]) {
	deps := getOrInit(m.Dependencies, key)

	for _, valKey := range deps.GetAll() {
		ok := r(valKey, m.Vals[valKey], m)
		if !ok {
			return
		}
	}
}

func (m *doubleMap[V]) RangeDependents(key any, r MutableRangeFunc[any, V]) {
	deps := getOrInit(m.Dependents, key)
	for _, valKey := range deps.GetAll() {
		ok := r(valKey, m.Vals[valKey], m)
		if !ok {
			return
		}
	}
}

func (m *doubleMap[V]) GetDependencies(key any) []any {
	d := getOrInit(m.Dependencies, key)
	return d.GetAll()
}

func (m *doubleMap[V]) GetDependents(key any) []any {
	d := getOrInit(m.Dependents, key)
	return d.GetAll()
}

func (m *doubleMap[V]) Remove(key any) {
	for _, v := range m.Dependents {
		v.Remove(key)
	}
	delete(m.Dependents, key)
	delete(m.Dependencies, key)
	delete(m.Vals, key)
}

func (m *doubleMap[V]) Add(key any, val V) {
	m.Vals[key] = val
}

func (m *doubleMap[V]) Keys() []any {
	keys := make([]any, len(m.Vals))
	i := 0
	for k := range m.Vals {
		keys[i] = k
		i++
	}
	return keys
}

func getOrInit(ma map[any]Set[any], key any) Set[any] {
	out := ma[key]
	if out == nil {
		out = NewSet[any]()
		ma[key] = out
	}
	return out
}

func (m *doubleMap[V]) Get(key any) V {
	return m.Vals[key]
}

func (m *doubleMap[V]) Lookup(key any) (V, bool) {
	val, ok := m.Vals[key]
	return val, ok
}
