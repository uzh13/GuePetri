package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
)

type mockHandlerForPlace[T any, V comparable] struct {
	handleInErr  error
	handleOutErr error
	chooseToErr  error
	chooseToRes  *graph.Transition[T, V]
}

func (m *mockHandlerForPlace[T, V]) HandleIn(*graph.Place[T, V]) error {
	return m.handleInErr
}

func (m *mockHandlerForPlace[T, V]) HandleOut(*graph.Place[T, V]) error {
	return m.handleOutErr
}

func (m *mockHandlerForPlace[T, V]) ChooseTo(T) (*graph.Transition[T, V], error) {
	return m.chooseToRes, m.chooseToErr
}

func TestNewPlace(t *testing.T) {
	handler := &mockHandlerForPlace[string, string]{}
	place := graph.NewPlace[string, string]("place1", handler)

	assert.Equal(t, "place1", place.ID)
	assert.Equal(t, handler, place.Handler)
}

func TestAddTransition(t *testing.T) {
	handler := &mockHandlerForPlace[string, string]{}
	place := graph.NewPlace("place1", handler)
	transition := &graph.Transition[string, string]{ID: "2"}

	place.AddTransition(transition)

	to := place.GetTo()
	_, exists := to[transition.ID]
	assert.True(t, exists)
}
