package graph

type TransitionHandler[T any, V comparable] interface {
	Handle(*Place[T, V], T) (*Place[T, V], error)
}

type Transition[T any, V comparable] struct {
	ID      V `json:"id,omitempty"`
	Handler TransitionHandler[T, V]
	to      map[V]struct{}
}

func NewTransition[T any, V comparable](id V, handler TransitionHandler[T, V]) *Transition[T, V] {
	return &Transition[T, V]{
		ID:      id,
		Handler: handler,
		to:      make(map[V]struct{}),
	}
}

func (t *Transition[T, V]) AddTo(n *Place[T, V]) *Transition[T, V] {
	if t.to == nil {
		t.to = make(map[V]struct{})
	}

	t.to[n.ID] = struct{}{}

	return t
}

func (p *Transition[T, V]) GetTo() map[V]struct{} {
	return p.to
}
