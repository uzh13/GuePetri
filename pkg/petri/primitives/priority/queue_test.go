package priority_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/priority"
)

func TestPriorityQueue_PushPeekPop(t *testing.T) {
	q := priority.NewPriorityQueue[int, int]()

	obj1 := &graph.Petri[int, int]{}
	obj2 := &graph.Petri[int, int]{}
	q.Push(1, obj1)
	q.Push(2, obj2)

	peeked, priority, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, obj2, peeked)
	assert.Equal(t, 2, priority)

	popped, ok := q.Pop()
	assert.True(t, ok)
	assert.Equal(t, obj2, popped)

	popped, ok = q.Pop()
	assert.True(t, ok)
	assert.Equal(t, obj1, popped)

	_, ok = q.Pop()
	assert.False(t, ok)
}

func TestPriorityQueue_PopPriority(t *testing.T) {
	q := priority.NewPriorityQueue[string, int]()

	obj1 := &graph.Petri[string, int]{}
	obj2 := &graph.Petri[string, int]{}
	q.Push(3, obj1)
	q.Push(1, obj2)

	popped, ok := q.PopPriority(1)
	assert.True(t, ok)
	assert.Equal(t, obj2, popped)

	_, ok = q.PopPriority(1)
	assert.False(t, ok)

	popped, ok = q.PopPriority(3)
	assert.True(t, ok)
	assert.Equal(t, obj1, popped)
}

func TestPriorityQueue_EmptyPeekPop(t *testing.T) {
	q := priority.NewPriorityQueue[float64, int]()

	_, _, ok := q.Peek()
	assert.False(t, ok)

	_, ok = q.Pop()
	assert.False(t, ok)
}
