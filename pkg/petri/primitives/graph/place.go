package graph

type PlaceHandler[T any, V comparable] interface {
	HandleIn(*Place[T, V]) error
	HandleOut(*Place[T, V]) error
	ChooseTo(T) (*Transition[T, V], error)
}

type Place[T any, V comparable] struct {
	ID      V `json:"id,omitempty"`
	to      map[V]struct{}
	handler PlaceHandler[T, V]
}

func NewNode[T any, V comparable](id V, handler PlaceHandler[T, V]) *Place[T, V] {
	return &Place[T, V]{
		ID:      id,
		to:      make(map[V]struct{}),
		handler: handler,
	}
}

func (p *Place[T, V]) AddTransition(s *Transition[T, V]) *Place[T, V] {
	p.to[s.ID] = struct{}{}

	return p
}
