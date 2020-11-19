package gui

import "errors"

var (
	// StateKeyError StateKeyError
	StateKeyError = errors.New("State key not existed. ")
)

// State State
type State interface {
	Set(key string, val interface{}) error
	Get(Ket string) (interface{}, error)
}

// StateMap StateMap
type StateMap struct {
	state map[string]interface{}
}

// NewStateMap NewStateMap
func NewStateMap() *StateMap {
	return &StateMap{state: map[string]interface{}{}}
}

// Set Set
func (s *StateMap) Set(key string, val interface{}) error {
	s.state[key] = val
	return nil
}

// Get Get
func (s *StateMap) Get(key string) (interface{}, error) {
	val, ok := s.state[key]
	if !ok {
		return nil, StateKeyError
	}

	return val, nil
}

// TowHeadQueue TowHeadQueue
type TowHeadQueue interface {
	Pop() interface{}
	Peek() interface{}
	Tail() interface{}
	Push(interface{})
	PopTail() interface{}
	Len() int
	IsEmpty() bool
}

// Queue Queue
type Queue struct {
	arr    []interface{}
	length int
}

// NewQueue NewQueue
func NewQueue() *Queue {
	return &Queue{
		arr:    make([]interface{}, 0),
		length: 0,
	}
}

// Pop Pop
func (q *Queue) Pop() interface{} {
	if q.length == 0 {
		return nil
	}

	index := q.length - 1
	el := q.arr[index]
	q.arr = q.arr[:index]
	q.length--
	return el
}

// Peek Peek
func (q *Queue) Peek() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.arr[q.length-1]
}

// Tail Tail
func (q *Queue) Tail() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.arr[0]
}

// PopTail PopTail
func (q *Queue) PopTail() interface{} {
	if q.length == 0 {
		return nil
	}
	el := q.arr[0]
	q.arr = q.arr[1:]
	q.length--
	return el
}

// Push Push
func (q *Queue) Push(el interface{}) {
	q.length++
	q.arr = append(q.arr, el)
}

// Len Len
func (q *Queue) Len() int {
	return q.length
}

// IsEmpty IsEmpty
func (q *Queue) IsEmpty() bool {
	return q.length == 0
}
