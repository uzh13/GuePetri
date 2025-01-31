package queue

type Queue[T any] struct {
	Elements []*T `json:"elements"`
}

// NewQueue создает новый экземпляр очереди
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		Elements: []*T{},
	}
}

// Enqueue добавляет элемент в конец очереди
func (q *Queue[T]) Enqueue(value *T) {
	q.Elements = append(q.Elements, value)
}

// Push alias to Enqueue
func (q *Queue[T]) Push(value *T) {
	q.Enqueue(value)
}

// Dequeue удаляет и возвращает первый элемент из очереди
func (q *Queue[T]) Dequeue() (*T, bool) {
	if len(q.Elements) == 0 {
		return nil, false
	}

	val := q.Elements[0]
	q.Elements = q.Elements[1:]

	return val, true
}

// Pop alias to Dequeue
func (q *Queue[T]) Pop() (*T, bool) {
	return q.Dequeue()
}

// Peek возвращает первый элемент очереди без удаления
func (q *Queue[T]) Peek() (*T, bool) {
	if len(q.Elements) == 0 {
		return nil, false
	}
	return q.Elements[0], true
}

// IsEmpty проверяет, пуста ли очередь
func (q *Queue[T]) IsEmpty() bool {
	return len(q.Elements) == 0
}
