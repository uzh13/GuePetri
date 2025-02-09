package graph

import (
	"fmt"
)

type PetriHandler interface {
	HandleIn() error
	HandleOut() error
}

type Petri[T any, V comparable] struct {
	ID      V            `json:"id"`
	Start   *Place[T, V] `json:"start"`
	Finish  *Place[T, V] `json:"finish"`
	Current *Place[T, V] `json:"current"`
	Handler PetriHandler
}

func NewPetri[T any, V comparable](id V, handler PetriHandler) *Petri[T, V] {
	return &Petri[T, V]{
		ID:      id,
		Handler: handler,
	}
}

func (g *Petri[T, V]) SetStartPlace(n *Place[T, V]) *Petri[T, V] {
	g.Start = n

	return g
}

func (g *Petri[T, V]) SetFinishPlace(n *Place[T, V]) *Petri[T, V] {
	g.Finish = n

	return g
}

func (g *Petri[T, V]) SetCurrentPlace(n *Place[T, V]) *Petri[T, V] {
	g.Current = n

	return g
}

func (g *Petri[T, V]) StartGraph() error {
	err := g.Handler.HandleIn()
	if err != nil {
		return fmt.Errorf("starting graph %v: %w", g.ID, err)
	}

	g.Current = g.Start
	err = g.Current.Handler.HandleIn(nil)
	if err != nil {
		return fmt.Errorf("starting graph %v first place %v: %w", g.ID, g.Current.ID, err)
	}

	return nil
}

func (g *Petri[T, V]) IsOnStart() bool {
	return g.Start.ID == g.Current.ID
}

func (g *Petri[T, V]) FinishGraph() error {
	err := g.Current.Handler.HandleOut(nil)
	if err != nil {
		return fmt.Errorf("finishing graph %v last place %v: %w", g.ID, g.Current.ID, err)
	}

	err = g.Handler.HandleOut()
	if err != nil {
		return fmt.Errorf("finishing graph %v: %w", g.ID, err)
	}

	return nil
}

func (g *Petri[T, V]) IsOnFinish() bool {
	return g.Finish.ID == g.Current.ID
}

func (g *Petri[T, V]) Act(signal T) error {
	current := g.Current
	if current == nil {
		err := g.StartGraph()
		if err != nil {
			return fmt.Errorf("auto starting graph : %w", err)
		}
	}

	transition, err := g.Current.Handler.ChooseTo(signal)
	if err != nil {
		return fmt.Errorf("graph %v choosing transition %v: %w", g.ID, g.Current.ID, err)
	}

	_, ok := g.Current.to[transition.ID]
	if !ok {
		return fmt.Errorf("graph %v forbitten transition %v for place %v", g.ID, transition.ID, g.Current.ID)
	}

	nextPlace, err := transition.Handler.Handle(g.Current, signal)
	if err != nil {
		return fmt.Errorf("graph %v handling signal %v by transition %v: %w", g.ID, signal, transition.ID, err)
	}

	_, ok = transition.to[nextPlace.ID]
	if !ok {
		return fmt.Errorf("graph %v forbitten place %v calculated from transition %v", g.ID, nextPlace.ID, transition.ID)
	}

	err = g.Current.Handler.HandleOut(nextPlace)
	if err != nil {
		return fmt.Errorf("graph %v exiting place %v to %v: %w", g.ID, g.Current.ID, nextPlace.ID, err)
	}

	oldCurrent := g.Current
	g.Current = nextPlace

	err = g.Current.Handler.HandleIn(oldCurrent)
	if err != nil {
		return fmt.Errorf("graph %v enterfing %v from %v: %w", g.ID, g.Current.ID, oldCurrent.ID, err)
	}

	return nil
}
