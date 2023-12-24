// Package event .
package event

import "context"

// Event defines the event type.
type Event string

// Listerner defines the callback function signature.
type Listerner[T any] func(context.Context, T)

// Emitter emits an event with associated data.
type Emitter[T any] struct {
	listeners map[Event][]Listerner[T]
}

// NewEmitter creates a new emitter.
func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		listeners: make(map[Event][]Listerner[T]),
	}
}

// On registers a callback for an event.
// It returns the emitter to allow chaining.
func (e *Emitter[T]) On(ev Event, l ...Listerner[T]) *Emitter[T] {
	e.listeners[ev] = append(e.listeners[ev], l...)
	return e
}

// Emit emits an event with associated data.
func (e *Emitter[T]) Emit(ctx context.Context, ev Event, data T) {
	l := e.listeners[ev]
	for _, f := range l {
		go f(ctx, data)
	}
}
