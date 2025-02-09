package aggregate_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/uzh13/GuePetri/pkg/petri/aggregate"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/priority"
)

type buffer struct {
	current string
}

func (b *buffer) Add(s string) {
	b.current = b.current + s + "\n"
}

func (b *buffer) Result() string {
	return b.current
}

type placeHandler struct {
	next      *graph.Transition[string, string]
	buffer    *buffer
	graphName string
	placeName string
}

func (h *placeHandler) HandleIn(from *graph.Place[string, string]) error {
	if from == nil {
		h.buffer.Add(fmt.Sprintf("handle in to %s from NIL, graph %s", h.placeName, h.graphName))

		return nil
	}
	h.buffer.Add(fmt.Sprintf("handle in to %s from %s, graph %s", h.placeName, from.ID, h.graphName))

	return nil
}

func (h *placeHandler) HandleOut(to *graph.Place[string, string]) error {
	if to == nil {
		h.buffer.Add(fmt.Sprintf("handle out from %s to NIL, graph %s", h.placeName, h.graphName))
		return nil
	}

	h.buffer.Add(fmt.Sprintf("handle out from %s to %s, graph %s", h.placeName, to.ID, h.graphName))

	return nil
}

func (h *placeHandler) ChooseTo(signal string) (*graph.Transition[string, string], error) {
	h.buffer.Add(fmt.Sprintf("ChooseTo at %s, signal %s, graph %s", h.placeName, signal, h.graphName))

	return h.next, nil
}

type transitionHandler struct {
	next           *graph.Place[string, string]
	buffer         *buffer
	graphName      string
	transitionName string
}

func (h *transitionHandler) Handle(p *graph.Place[string, string], signal string) (*graph.Place[string, string], error) {
	h.buffer.Add(fmt.Sprintf("signal %s, transition %s, graph %s", signal, h.transitionName, h.graphName))

	return h.next, nil
}

type graphHandler struct {
	buffer    *buffer
	graphName string
}

func (h *graphHandler) HandleIn() error {
	h.buffer.Add(fmt.Sprintf("handle in graph, graph %s", h.graphName))

	return nil
}

func (h *graphHandler) HandleOut() error {
	h.buffer.Add(fmt.Sprintf("handle out graph %s", h.graphName))

	return nil
}

func TestPetriQueue_Act_OneLevel(t *testing.T) {
	b := buffer{current: "\n"}
	expected := `
handle in graph, graph graph1
handle in to start from NIL, graph graph1
ChooseTo at start, signal sig0, graph graph1
signal sig0, transition middle, graph graph1
handle out from start to finish, graph graph1
handle in to finish from start, graph graph1
handle out from finish to NIL, graph graph1
handle out graph graph1
handle in graph, graph graph2
handle in to start from NIL, graph graph2
ChooseTo at start, signal 0, graph graph2
signal 0, transition start_to_finish, graph graph2
handle out from start to middle, graph graph2
handle in to middle from start, graph graph2
ChooseTo at middle, signal sig1, graph graph2
signal sig1, transition middle_to_finish, graph graph2
handle out from middle to finish, graph graph2
handle in to finish from middle, graph graph2
handle out from finish to NIL, graph graph2
handle out graph graph2
`
	target := aggregate.NewPetriQueue[string, string](priority.NewPriorityQueue[string, string](), "0")
	err := target.AddGraph(0, makeGraph1(&b))
	if err != nil {
		t.Errorf("failed to add first graph")
	}

	err = target.AddGraph(0, makeGraph2(&b))
	if err != nil {
		t.Errorf("failed to add second graph")
	}

	for i := 0; i < 3; i++ {
		err = target.Act("sig" + strconv.Itoa(i))
		if err != nil {
			t.Errorf("iteration %d, %v", i, err)
		}
	}

	if expected != b.Result() {
		t.Errorf("expected: %s, got: %s", expected, b.Result())
	}
}

func TestPetriQueue_Act_TwoLevels(t *testing.T) {
	b := buffer{current: "\n"}
	expected := `
handle in graph, graph graph1
handle in to start from NIL, graph graph1
handle in graph, graph graph2
handle in to start from NIL, graph graph2
ChooseTo at start, signal sig0, graph graph2
signal sig0, transition start_to_finish, graph graph2
handle out from start to middle, graph graph2
handle in to middle from start, graph graph2
ChooseTo at middle, signal sig1, graph graph2
signal sig1, transition middle_to_finish, graph graph2
handle out from middle to finish, graph graph2
handle in to finish from middle, graph graph2
handle out from finish to NIL, graph graph2
handle out graph graph2
ChooseTo at start, signal 0, graph graph1
signal 0, transition middle, graph graph1
handle out from start to finish, graph graph1
handle in to finish from start, graph graph1
handle out from finish to NIL, graph graph1
handle out graph graph1
`
	target := aggregate.NewPetriQueue[string, string](priority.NewPriorityQueue[string, string](), "0")

	err := target.AddGraph(0, makeGraph1(&b))
	if err != nil {
		t.Errorf("failed to add first graph")
	}

	err = target.AddGraph(1, makeGraph2(&b))
	if err != nil {
		t.Errorf("failed to add second graph")
	}

	for i := 0; i < 3; i++ {
		err := target.Act("sig" + strconv.Itoa(i))
		if err != nil {
			t.Errorf("iteration %d, %v", i, err)
		}
	}

	if expected != b.Result() {
		t.Errorf("expected: %s, got: %s", expected, b.Result())
	}
}

func makeGraph1(b *buffer) *graph.Petri[string, string] {
	const name = "graph1"

	placeFinishHandler := placeHandler{
		next:      nil,
		buffer:    b,
		graphName: name,
		placeName: "finish",
	}
	placeFinish := graph.NewPlace[string, string]("finish", &placeFinishHandler)

	trHandler := transitionHandler{
		next:           placeFinish,
		buffer:         b,
		graphName:      name,
		transitionName: "middle",
	}
	transition := graph.NewTransition[string, string]("middle", &trHandler).AddTo(placeFinish)

	placeStartHandler := placeHandler{
		next:      transition,
		buffer:    b,
		graphName: name,
		placeName: "start",
	}
	placeStart := graph.NewPlace[string, string]("start", &placeStartHandler).AddTransition(transition)

	grHandler := graphHandler{
		buffer:    b,
		graphName: name,
	}

	return graph.NewPetri[string, string](name, &grHandler).
		SetStartPlace(placeStart).
		SetFinishPlace(placeFinish)
}

func makeGraph2(b *buffer) *graph.Petri[string, string] {
	const name = "graph2"

	placeFinishHandler := placeHandler{
		next:      nil,
		buffer:    b,
		graphName: name,
		placeName: "finish",
	}
	placeFinish := graph.NewPlace[string, string]("finish", &placeFinishHandler)

	trHandler2 := transitionHandler{
		next:           placeFinish,
		buffer:         b,
		graphName:      name,
		transitionName: "middle_to_finish",
	}
	transition2 := graph.NewTransition[string, string]("middle_to_finish", &trHandler2).
		AddTo(placeFinish)

	placeMiddleHandler := placeHandler{
		next:      transition2,
		buffer:    b,
		graphName: name,
		placeName: "middle",
	}
	placeMiddle := graph.NewPlace[string, string]("middle", &placeMiddleHandler).
		AddTransition(transition2)

	trHandler := transitionHandler{
		next:           placeMiddle,
		buffer:         b,
		graphName:      name,
		transitionName: "start_to_finish",
	}
	transition := graph.NewTransition[string, string]("start_to_finish", &trHandler).
		AddTo(placeMiddle)

	placeStartHandler := placeHandler{
		next:      transition,
		buffer:    b,
		graphName: name,
		placeName: "start",
	}
	placeStart := graph.NewPlace[string, string]("start", &placeStartHandler).
		AddTransition(transition)

	grHandler := graphHandler{
		buffer:    b,
		graphName: name,
	}

	return graph.NewPetri[string, string](name, &grHandler).
		SetStartPlace(placeStart).
		SetFinishPlace(placeFinish)
}
