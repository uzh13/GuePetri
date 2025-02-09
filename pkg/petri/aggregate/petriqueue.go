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

func (p *PetriQueue[T, V]) AddGraph(level int, graph *graph.Petri[T, V]) error {
	if p.queue == nil {
		p.queue = priority.NewPriorityQueue[T, V]()
	}

	size := p.queue.Len()
	currentLevel := p.queue.GetMaxPriority()

	p.queue.Push(level, graph)

	if graph.Current != nil {
		return nil
	}

	if size != 0 && level <= currentLevel {
		return nil
	}

	err := graph.StartGraph()
	if err != nil {
		return fmt.Errorf("add graph and start: %w", err)
	}

	return nil
}

func (p *PetriQueue[T, V]) GetQueue() *priority.Queue[T, V] {
	return p.queue
}

func (p *PetriQueue[T, V]) Act(signal T) error {
	if p.queue == nil {
		p.queue = priority.NewPriorityQueue[T, V]()
	}

	if len(p.queue.GrQu) == 0 {
		return nil
	}

	current, priorityLevel, ok := p.queue.Peek()
	if !ok {
		return errors.New("unable to detect peek")
	}

	err := current.Act(signal)
	if err != nil {
		return fmt.Errorf("unable to act priority %v, graph %v, signal %v: %w", priorityLevel, current, signal, err)
	}

	if !current.IsOnFinish() {
		return nil
	}

	err = current.FinishGraph()
	if err != nil {
		return fmt.Errorf("unable to finish graph %v, priority %v, signal %v: %w", current, priorityLevel, signal, err)
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
