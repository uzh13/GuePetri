package graph_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
)

type mockHandler struct {
	inErr  error
	outErr error
}

func (m *mockHandler) HandleIn() error {
	return m.inErr
}

func (m *mockHandler) HandleOut() error {
	return m.outErr
}

type mockPlaceHandle struct {
	inErr     error
	outErr    error
	chooseErr error
	choose    *graph.Transition[int, string]
}

func (m *mockPlaceHandle) HandleIn(place *graph.Place[int, string]) error {
	return m.inErr
}

func (m *mockPlaceHandle) HandleOut(place *graph.Place[int, string]) error {
	return m.outErr
}

func (m *mockPlaceHandle) ChooseTo(signal int) (*graph.Transition[int, string], error) {
	if m.choose == nil {
		return nil, m.chooseErr
	}

	return m.choose, nil
}

type mocktransitionHandler struct {
	result *graph.Place[int, string]
}

func (m *mocktransitionHandler) Handle(p *graph.Place[int, string], t int) (*graph.Place[int, string], error) {
	return m.result, nil
}

func TestPetri_StartGraph(t *testing.T) {
	handler := &mockHandler{}
	placeHandler := &mockPlaceHandle{}
	startPlace := &graph.Place[int, string]{ID: "start", Handler: placeHandler}
	finishPlace := &graph.Place[int, string]{ID: "finish", Handler: placeHandler}

	petri := &graph.Petri[int, string]{
		ID:      "testGraph",
		Start:   startPlace,
		Finish:  finishPlace,
		Handler: handler,
	}

	err := petri.StartGraph()
	assert.NoError(t, err)
	assert.Equal(t, startPlace, petri.Current)
}

func TestPetri_FinishGraph(t *testing.T) {
	handler := &mockHandler{}
	placeHandler := &mockPlaceHandle{}
	finishPlace := &graph.Place[int, string]{ID: "finish", Handler: placeHandler}

	petri := &graph.Petri[int, string]{
		ID:      "testGraph",
		Finish:  finishPlace,
		Current: finishPlace,
		Handler: handler,
	}

	err := petri.FinishGraph()
	assert.NoError(t, err)
}

func TestPetri_IsOnStartAndFinish(t *testing.T) {
	startPlace := &graph.Place[int, string]{ID: "start"}
	finishPlace := &graph.Place[int, string]{ID: "finish"}

	petri := &graph.Petri[int, string]{
		ID:      "testGraph",
		Start:   startPlace,
		Finish:  finishPlace,
		Current: startPlace,
	}

	assert.True(t, petri.IsOnStart())
	assert.False(t, petri.IsOnFinish())

	petri.Current = finishPlace
	assert.False(t, petri.IsOnStart())
	assert.True(t, petri.IsOnFinish())
}

func TestPetri_StartGraph_ErrorHandling(t *testing.T) {
	handler := &mockHandler{inErr: errors.New("handler error")}
	placeHandler := &mockPlaceHandle{}
	startPlace := &graph.Place[int, string]{ID: "start", Handler: placeHandler}
	finishPlace := &graph.Place[int, string]{ID: "finish", Handler: placeHandler}

	petri := &graph.Petri[int, string]{
		ID:      "testGraph",
		Start:   startPlace,
		Finish:  finishPlace,
		Handler: handler,
	}

	err := petri.StartGraph()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "handler error")
}

func TestPetri_Act(t *testing.T) {
	placeHandlerStart := &mockPlaceHandle{}
	startPlace := graph.NewPlace[int, string]("start", placeHandlerStart)

	placeHandlerNext := &mockPlaceHandle{}
	nextPlace := graph.NewPlace[int, string]("next", placeHandlerNext)

	transition := graph.NewTransition[int, string](
		"testTransition",
		&mocktransitionHandler{
			result: nextPlace,
		},
	)

	placeHandlerStart.choose = transition

	transition.AddTo(nextPlace)
	startPlace.AddTransition(transition)

	petri := &graph.Petri[int, string]{
		ID:      "testGraph",
		Start:   startPlace,
		Current: startPlace,
		Finish:  nextPlace,
		Handler: &mockHandler{},
	}

	err := petri.Act(1)
	assert.NoError(t, err)
	assert.Equal(t, "next", petri.Current.ID)
}
