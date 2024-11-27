package types

type MapSlice[T comparable] struct {
	mapData map[T]bool
	slice   []T
}

func NewSliceMap[T comparable]() *MapSlice[T] {
	return &MapSlice[T]{
		mapData: make(map[T]bool),
		slice:   make([]T, 0, MaxNumDefault),
	}
}

func (s *MapSlice[T]) Each() []T {
	return s.slice
}

func (s *MapSlice[T]) Exists(id T) bool {
	_, ok := s.mapData[id]
	return ok
}

func (s *MapSlice[T]) Len() int {
	return len(s.slice)
}

func (s *MapSlice[T]) Push(id T) {
	if s.Exists(id) {
		return
	}
	s.mapData[id] = true
	s.slice = append(s.slice, id)
}
func (s *MapSlice[T]) PushN(ids []T) {
	for _, id := range ids {
		if s.Exists(id) {
			continue
		}
		s.mapData[id] = true
		s.slice = append(s.slice, id)
	}
}

func (s *MapSlice[T]) Get(index int) T {
	if index > len(s.slice)-1 {
		var v T
		return v
	}

	element := s.slice[index]
	return element
}

func (s *MapSlice[T]) Pop() T {
	if len(s.slice) > 0 {
		return s.slice[0]
	}
	var v T
	return v
}

func (s *MapSlice[T]) Remove(id T) {
	delete(s.mapData, id)
	for i, v := range s.slice {
		if v == id {
			s.slice = append(s.slice[:i], s.slice[i+1:]...)
			break
		}
	}
}
