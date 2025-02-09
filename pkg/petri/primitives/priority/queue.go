package priority

import (
	"math"
	"sync"

	"github.com/uzh13/GuePetri/pkg/petri/primitives/graph"
	"github.com/uzh13/GuePetri/pkg/petri/primitives/queue"
)

type Queue[T any, V comparable] struct {
	GrQu        map[int]*queue.Queue[graph.Petri[T, V]] `json:"objects"`
	maxPriority int
	mu          *sync.Mutex
}

func NewPriorityQueue[T any, V comparable]() *Queue[T, V] {
	return &Queue[T, V]{
		GrQu: make(map[int]*queue.Queue[graph.Petri[T, V]]),
		mu:   &sync.Mutex{},
	}
}

// Peek Читаем актуальный объект, без изменения состояния
func (p *Queue[T, V]) Peek() (*graph.Petri[T, V], int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.GrQu) == 0 {
		return nil, 0, false
	}

	q, ok := p.GrQu[p.maxPriority]
	if !ok {
		return nil, 0, false
	}

	result, ok := q.Peek()
	if !ok {
		return nil, 0, false
	}

	return result, p.maxPriority, true
}

// Push добавить объект
func (p *Queue[T, V]) Push(priority int, obj *graph.Petri[T, V]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check for integer overflow or invalid input
	if priority > math.MaxInt {
		panic("invalid input: index out of range")
	}

	st, ok := p.GrQu[priority]
	if !ok {
		st = queue.NewQueue[graph.Petri[T, V]]()
		p.GrQu[priority] = st

		if priority > p.maxPriority {
			p.maxPriority = priority
		}
	}

	st.Push(obj)
}

// PopPriority выдёргивает актуальный элемент с определённого уровня приоритета
func (p *Queue[T, V]) PopPriority(priority int) (*graph.Petri[T, V], bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	st, ok := p.GrQu[priority]
	if !ok {
		return nil, false
	}

	obj, ok := st.Pop()
	if !ok {
		delete(p.GrQu, priority)
		if priority == p.maxPriority {
			p.maxPriority = calculateMax(p.GrQu)
		}

		return nil, false
	}

	if len(st.Elements) == 0 {
		delete(p.GrQu, priority)
		if priority == p.maxPriority {
			p.maxPriority = calculateMax(p.GrQu)
		}
	}

	return obj, true
}

func (p *Queue[T, V]) Pop() (*graph.Petri[T, V], bool) {
	return p.PopPriority(p.maxPriority)
}

func (p *Queue[T, V]) GetMaxPriority() int {
	return p.maxPriority
}

func (p *Queue[T, V]) Len() int {
	return len(p.GrQu)
}

func calculateMax[T any](m map[int]*queue.Queue[T]) int {
	if len(m) == 0 {
		return 0
	}

	var maxKey int

	started := false
	for k := range m {
		if !started {
			started = true
			maxKey = k
			continue
		}

		if maxKey < k {
			maxKey = k
		}
	}

	return maxKey
}
