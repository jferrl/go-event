// Package event .
package event

import "context"

// Event defines the event type.
type Event string

// Listerner defines the callback function signature.
type Listerner[T any] func(context.Context, T)

// EmitterOption configures the emitter.
type EmitterOption[T any] func(*Emitter[T])

// SyncEmitter configures the emitter to be synchronous.
// It means that listeners are executed sequentially.
// Default: false.
func SyncEmitter[T any](e *Emitter[T]) {
	e.sync = true
}

// Emitter emits an event with associated data.
type Emitter[T any] struct {
	listeners map[Event][]Listerner[T]
	sync      bool
}

// NewEmitter creates a new emitter.
// It accepts a list of options to configure the emitter.
// By default, the emitter is asynchronous. It means that each listener is
// executed in a separate goroutine.
func NewEmitter[T any](opts ...EmitterOption[T]) *Emitter[T] {
	e := &Emitter[T]{
		listeners: make(map[Event][]Listerner[T]),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// On registers a callback for an event.
// It returns the emitter to allow chaining.
func (e *Emitter[T]) On(ev Event, l ...Listerner[T]) *Emitter[T] {
	e.listeners[ev] = append(e.listeners[ev], l...)
	return e
}

// Emit emits an event with associated data.
// It does not block the caller.
func (e *Emitter[T]) Emit(ctx context.Context, ev Event, data T) {
	go func(l []Listerner[T]) {
		for _, f := range l {
			if e.sync {
				f(ctx, data)
			} else {
				go f(ctx, data)
			}
		}
	}(e.listeners[ev])
}
