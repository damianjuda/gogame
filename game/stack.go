package game

import (
    "sync"
)

// Item the type of the stack

// ItemStack the stack of Items
type ItemStack struct {
    items []Step
    lock  sync.RWMutex
}

// New creates a new ItemStack
func (s *ItemStack) New() *ItemStack {
    s.items = []Step{}
    return s
}

// Push adds an Item to the top of the stack
func (s *ItemStack) Push(t Step) {
    s.lock.Lock()
    s.items = append(s.items, t)
    s.lock.Unlock()
}

// Pop removes an Item from the top of the stack
func (s *ItemStack) Pop() Step {
	s.lock.Lock()
	if s.Len() > 0 {
		item := s.items[len(s.items)-1]
		s.items = s.items[0 : len(s.items)-1]
		s.lock.Unlock()
    	return item
	} else {
		s.lock.Unlock()
		return nil
	}
}

func (s *ItemStack) Len() int {
	return len(s.items)
}