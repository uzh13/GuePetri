package graph

type TransitionHandler[T any, V comparable] interface {
	Handle(*Place[T, V], T) (*Place[T, V], error)
}

type Transition[T any, V comparable] struct {
	ID      V `json:"id,omitempty"`
	to      map[V]struct{}
	handler TransitionHandler[T, V]
}

func NewTransition[T any, V comparable](id V, handler TransitionHandler[T, V]) *Transition[T, V] {
	return &Transition[T, V]{
		ID:      id,
		to:      make(map[V]struct{}),
		handler: handler,
	}
}

func (t *Transition[T, V]) AddTo(n *Place[T, V]) *Transition[T, V] {
	t.to[n.ID] = struct{}{}

	return t
}
