package aggregate

import (
	"fmt"

	"github.com/uzh13/GuePetri/pkg/petri/primitives/priority"
)

type Storage[T any, V comparable, U comparable] interface {
	Get(U) (*priority.Queue[T, V], error)
}

type Builder[T any, V comparable, U comparable] struct {
	ID      U
	Storage Storage[T, V, U]
	petriQ  *priority.Queue[T, V]
}

func NewBuilder[T any, V comparable, U comparable](id U, storage Storage[T, V, U]) *Builder[T, V, U] {
	return &Builder[T, V, U]{
		ID:      id,
		Storage: storage,
	}
}

func (b *Builder[T, V, U]) LoadState() error {
	stored, err := b.Storage.Get(b.ID)
	if err != nil {
		return fmt.Errorf("unable to load state: %w", err)
	}

	if stored == nil {
		b.petriQ = priority.NewPriorityQueue[T, V]()
	} else {
		b.petriQ = stored
	}

	return nil
}

func (b *Builder[T, V, U]) Build(zeroSignal T) *PetriQueue[T, V] {
	return NewPetriQueue(b.petriQ, zeroSignal)
}
