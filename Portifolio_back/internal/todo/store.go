package todo

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("Not Found")

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	next  int
	items map[int]Todo
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[int]Todo), next: 1}
}

func (s *MemoryStore) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Todo, 0, len(s.items))
	for _, t := range s.items {
		out = append(out, t)
	}
	return out
}

func (s *MemoryStore) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	t := Todo{ID: id, Title: title}
	s.items[id] = t
	return t
}

func (s *MemoryStore) Get(id int) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) Update(id int, title string, done bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	t := Todo{ID: id, Title: title, Done: done}
	s.items[id] = t
	return t, nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
