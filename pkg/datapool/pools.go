package datapool

type (
	Pool[K comparable, T any] struct {
		values map[K]T
	}
)

func NewPool[K comparable, T any]() Pool[K, T] {
	return Pool[K, T]{
		values: map[K]T{},
	}
}

func (p *Pool[K, T]) Add(key K, value T) {
	p.values[key] = value
}

func (p *Pool[K, T]) Has(key K) bool {
	_, ok := p.values[key]
	return ok
}

func (p *Pool[K, T]) Get(key K) (T, bool) {
	found, ok := p.values[key]
	return found, ok
}
