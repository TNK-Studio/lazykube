package gui

type State interface {
	Set(key string, val interface{}) error
	Get(Ket string) (interface{}, error)
}

type StateMap struct {
	state map[string]interface{}
}

func (s *StateMap) Set(key string, val interface{}) error {
	s.state[key] = val
	return nil
}

func (s *StateMap) Get(key string) (interface{}, error) {
	return s.state[key], nil
}
