package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
)

type mockTransitionHandler[T any, V comparable] struct {
	handleErr error
	handleRes *graph.Place[T, V]
}

func (m *mockTransitionHandler[T, V]) Handle(p *graph.Place[T, V], t T) (*graph.Place[T, V], error) {
	return m.handleRes, m.handleErr
}

func TestNewTransition(t *testing.T) {
	handler := &mockTransitionHandler[string, int]{}
	transition := graph.NewTransition(1, handler)

	assert.Equal(t, 1, transition.ID)
	assert.Equal(t, handler, transition.Handler)
}

func TestAddTo(t *testing.T) {
	handler := &mockTransitionHandler[string, string]{}
	transition := graph.NewTransition("1", handler)
	place := graph.NewPlace[string, string]("place1", nil)

	transition.AddTo(place)

	to := transition.GetTo()
	_, exists := to[place.ID]
	assert.True(t, exists)
}
