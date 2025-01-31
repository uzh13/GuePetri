package aggregate

import (
	"errors"
	"fmt"

	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/priority"
)

type PetriQueue[T any, V comparable] struct {
	queue      *priority.Queue[T, V]
	zeroSignal T
}

func NewPetriQueue[T any, V comparable](q *priority.Queue[T, V], zero T) *PetriQueue[T, V] {
	return &PetriQueue[T, V]{
		queue:      q,
		zeroSignal: zero,
	}
}

func (p *PetriQueue[T, V]) AddGraph(priority int, graph *graph.Petri[T, V]) {
	if p.queue == nil {
		p.queue = priority.NewPriorityQueue[T, V]()
	}

	p.queue.Push(priority, graph)
}

func (p *PetriQueue[T, V]) GetQueue() *priority.Queue[T, V] {
	return p.queue
}

func (p *PetriQueue[T, V]) Act(signal T) error {
	if p.queue == nil {
		p.queue = priority.NewPriorityQueue[T, V]()
	}

	current, priorityLevel, ok := p.queue.Peek()
	if !ok {
		return errors.New("unable to detect peek")
	}

	err := current.Act(signal)
	if err != nil {
		return fmt.Errorf("unable to act priority %v, graph %s, signal %v: %w", priorityLevel, current, signal, err)
	}

	if !current.IsOnFinish() {
		return nil
	}

	err = current.FinishGraph()
	if err != nil {
		return fmt.Errorf("unable to finish graph %s, priority %v, signal %v: %w", current, priorityLevel, signal, err)
	}

	graphToRemove, ok := p.queue.PopPriority(priorityLevel)
	if !ok {
		return errors.New(fmt.Sprintf("unable to get graph to delete, level %d, signal %v", priorityLevel, signal))
	}

	if graphToRemove.ID != current.ID {
		return errors.New(
			fmt.Sprintf(
				"ID of graphs are not identical, level %d, graph %v, to delete %v",
				priorityLevel,
				graphToRemove.ID,
				current.ID,
			),
		)
	}

	err = p.Act(p.zeroSignal)
	if err != nil {
		return fmt.Errorf("act after current delete with zero signal: %w", err)
	}

	return nil
}
