package app

import "fmt"

type commands struct {
	handlers map[string]func(*state, command) error
}

func newCommands() *commands {
	return &commands{
		handlers: make(map[string]func(*state, command) error),
	}
}

func (registry *commands) register(name string, handler func(*state, command) error) {
	if registry.handlers == nil {
		registry.handlers = make(map[string]func(*state, command) error)
	}
	registry.handlers[name] = handler
}

func (registry *commands) run(st *state, cmd command) error {
	handler, exists := registry.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("command %q not found", cmd.name)
	}

	return handler(st, cmd)
}
