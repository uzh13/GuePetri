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

func (g *Petri[T, V]) SetStartNode(n *Place[T, V]) *Petri[T, V] {
	g.Start = n

	return g
}

func (g *Petri[T, V]) SetFinishNode(n *Place[T, V]) *Petri[T, V] {
	g.Finish = n

	return g
}

func (g *Petri[T, V]) SetCurrentNode(n *Place[T, V]) *Petri[T, V] {
	g.Current = n

	return g
}

func (g *Petri[T, V]) StartGraph() error {
	err := g.Handler.HandleIn()
	if err != nil {
		return fmt.Errorf("starting graph %s: %w", g.ID, err)
	}

	g.Current = g.Start
	err = g.Current.handler.HandleIn(nil)
	if err != nil {
		return fmt.Errorf("starting graph %s first node %s: %w", g.ID, g.Current.ID, err)
	}

	return nil
}

func (g *Petri[T, V]) IsOnStart() bool {
	return g.Start.ID == g.Current.ID
}

func (g *Petri[T, V]) FinishGraph() error {
	err := g.Current.handler.HandleOut(nil)
	if err != nil {
		return fmt.Errorf("finishing graph %s last node %s: %w", g.ID, g.Current.ID, err)
	}

	err = g.Handler.HandleOut()
	if err != nil {
		return fmt.Errorf("finishing graph %s: %w", g.ID, err)
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

	nextPeer, err := g.Current.handler.ChooseTo(signal)
	if err != nil {
		return fmt.Errorf("graph %s choosing peer %s: %w", g.ID, g.Current.ID, err)
	}

	_, ok := g.Current.to[nextPeer.ID]
	if !ok {
		return fmt.Errorf("graph %s forbitten peer %s for node %s", g.ID, nextPeer.ID, g.Current.ID)
	}

	nextNode, err := nextPeer.handler.Handle(g.Current, signal)
	if err != nil {
		return fmt.Errorf("graph %s handling signal %s by peer %s: %w", g.ID, signal, nextPeer.ID, err)
	}

	_, ok = nextPeer.to[nextNode.ID]
	if !ok {
		return fmt.Errorf("graph %s forbitten node %s calculated from peer %s", g.ID, nextNode.ID, nextPeer.ID)
	}

	err = g.Current.handler.HandleOut(nextNode)
	if err != nil {
		return fmt.Errorf("graph %s exiting node %s to %s: %w", g.ID, g.Current.ID, nextNode.ID, err)
	}

	oldCurrent := g.Current
	g.Current = nextNode

	err = g.Current.handler.HandleIn(oldCurrent)
	if err != nil {
		return fmt.Errorf("graph %s enterfing %s from %s: %w", g.ID, g.Current, oldCurrent, err)
	}

	return nil
}
