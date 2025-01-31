package queue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/queue" // Замените на правильный путь
)

func TestQueue_EnqueueDequeue(t *testing.T) {
	q := queue.NewQueue[int]()

	val1, val2 := 10, 20
	q.Enqueue(&val1)
	q.Enqueue(&val2)

	deqVal, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, val1, *deqVal)

	deqVal, ok = q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, val2, *deqVal)

	_, ok = q.Dequeue()
	assert.False(t, ok)
}

func TestQueue_PushPop(t *testing.T) {
	q := queue.NewQueue[int]()

	val := 42
	q.Push(&val)

	popped, ok := q.Pop()
	assert.True(t, ok)
	assert.Equal(t, val, *popped)

	_, ok = q.Pop()
	assert.False(t, ok)
}

func TestQueue_Peek(t *testing.T) {
	q := queue.NewQueue[string]()

	val := "hello"
	q.Enqueue(&val)

	peeked, ok := q.Peek()
	assert.True(t, ok)
	assert.Equal(t, val, *peeked)

	deqVal, ok := q.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, val, *deqVal)

	_, ok = q.Peek()
	assert.False(t, ok)
}

func TestQueue_IsEmpty(t *testing.T) {
	q := queue.NewQueue[float64]()
	assert.True(t, q.IsEmpty())

	val := 3.14
	q.Enqueue(&val)
	assert.False(t, q.IsEmpty())

	_, _ = q.Dequeue()
	assert.True(t, q.IsEmpty())
}
