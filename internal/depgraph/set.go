package depgraph

type Set[K any] interface {
	Add(key K)
	Remove(key K)
	GetAll() []K
}

func NewSet[K any]() Set[K] {
	return &set[K]{
		Map: map[any]bool{},
	}
}

type set[K any] struct {
	Map map[any]bool
}

func (s *set[K]) GetAll() []K {
	out := make([]K, len(s.Map))
	i := 0
	for k := range s.Map {
		out[i] = k.(K)
		i++
	}
	return out
}

func (s *set[K]) Add(key K) {
	s.Map[key] = true
}

func (s *set[K]) Remove(key K) {
	delete(s.Map, key)
}
