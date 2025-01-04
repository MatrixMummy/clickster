package main

import (
	"fmt"
	"sync"
)

// StateManager is responsible for managing the state of the application
type StateManager struct {
	mu    sync.Mutex
	state map[string]interface{}
}

// NewStateManager creates a new StateManager instance
func NewStateManager() *StateManager {
	fmt.Println("Creating new StateManager instance")
	return &StateManager{
		state: make(map[string]interface{}),
	}
}

// Set sets a value in the state
func (sm *StateManager) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state[key] = value
}

// Get retrieves a value from the state
func (sm *StateManager) Get(key string) (interface{}, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	value, exists := sm.state[key]
	return value, exists
}

// // NewEngine creates a new Engine instance
// func NewEngine() *Engine {
// 	fmt.Println("Creating new Engine instance")
// 	return &Engine{}
// }
