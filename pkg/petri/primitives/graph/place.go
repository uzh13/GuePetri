package graph

type PlaceHandler[T any, V comparable] interface {
	HandleIn(*Place[T, V]) error
	HandleOut(*Place[T, V]) error
	ChooseTo(T) (*Transition[T, V], error)
}

type Place[T any, V comparable] struct {
	ID      V `json:"id,omitempty"`
	Handler PlaceHandler[T, V]
	to      map[V]struct{}
}

func NewPlace[T any, V comparable](id V, handler PlaceHandler[T, V]) *Place[T, V] {
	return &Place[T, V]{
		ID:      id,
		Handler: handler,
		to:      make(map[V]struct{}),
	}
}

func (p *Place[T, V]) AddTransition(s *Transition[T, V]) *Place[T, V] {
	if p.to == nil {
		p.to = make(map[V]struct{})
	}

	p.to[s.ID] = struct{}{}

	return p
}
