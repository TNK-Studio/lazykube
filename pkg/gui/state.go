package gui

import "errors"

var (
	StateKeyError = errors.New("State key not existed. ")
)

type State interface {
	Set(key string, val interface{}) error
	Get(Ket string) (interface{}, error)
}

type StateMap struct {
	state map[string]interface{}
}

func NewStateMap() *StateMap {
	return &StateMap{state: map[string]interface{}{}}
}

func (s *StateMap) Set(key string, val interface{}) error {
	s.state[key] = val
	return nil
}

func (s *StateMap) Get(key string) (interface{}, error) {
	val, ok := s.state[key]
	if !ok {
		return nil, StateKeyError
	}

	return val, nil
}

type TowHeadQueue interface {
	Pop() interface{}
	Peek() interface{}
	Tail() interface{}
	Push(interface{})
	PopTail() interface{}
	Len() int
	IsEmpty() bool
}

type Queue struct {
	arr    []interface{}
	length int
}

func NewQueue() *Queue {
	return &Queue{
		arr:    make([]interface{}, 0),
		length: 0,
	}
}

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

func (q *Queue) Peek() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.arr[q.length-1]
}

func (q *Queue) Tail() interface{} {
	if q.length == 0 {
		return nil
	}
	return q.arr[0]
}

func (q *Queue) PopTail() interface{} {
	if q.length == 0 {
		return nil
	}
	el := q.arr[0]
	q.arr = q.arr[1:]
	q.length--
	return el
}

func (q *Queue) Push(el interface{}) {
	q.length++
	q.arr = append(q.arr, el)
}

func (q *Queue) Len() int {
	return q.length
}

func (q *Queue) IsEmpty() bool {
	return q.length == 0
}
