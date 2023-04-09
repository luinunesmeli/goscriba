package datapool

type (
	Pool[T comparable] struct {
		values []T
	}
)

func (p *Pool[T]) Add(value ...T) {
	p.values = append(p.values, value...)
}

func (p *Pool[T]) Has(value T) bool {
	for _, val := range p.values {
		if val == value {
			return true
		}
	}
	return false
}
