# go-event

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/jferrl/go-event)
[![Test Status](https://github.com/jferrl/go-event/workflows/tests/badge.svg)](https://github.com/jferrl/go-event/actions?query=workflow%3Atests)
[![codecov](https://codecov.io/gh/jferrl/go-event/branch/main/graph/badge.svg?token=68I4BZF235)](https://codecov.io/gh/jferrl/go-event)
[![Go Report Card](https://goreportcard.com/badge/github.com/jferrl/go-event)](https://goreportcard.com/report/github.com/jferrl/go-event)

Go simple (zero deps) library for event handling. The main goal of this library is to provide a simple and easy to use event handling system with a minimal footprint within a Go application. This library is inspired by Node.js [EventEmitter](https://nodejs.org/api/events.html#events_class_eventemitter).

## Usage

go-event is compatible with modern Go releases.
This pkg use Go generics, so you need to use Go 1.18 or later.

```go
package main

import (
 "context"
 "log"
 "os"
 "time"

 "github.com/jferrl/go-event"
)

const (
 // UserCreated .
 UserCreated event.Event = "user.created"
 // UserDeleted .
 UserDeleted event.Event = "user.deleted"
 // UserUpdated .
 UserUpdated event.Event = "user.updated"
)

// UserEvent .
type UserEvent struct {
 ID string
}

func main() {
 logger := log.New(os.Stdout, "go-event: ", log.LstdFlags)

 ctx := context.Background()

 emitter := event.NewEmitter[UserEvent]()

 emitter.
  On(UserCreated, func(ctx context.Context, data UserEvent) {
   // handle user created event
   logger.Printf("user created: %s", data.ID)
  }).
  On(UserDeleted, func(ctx context.Context, data UserEvent) {
   // handle user deleted event
   logger.Printf("user deleted: %s", data.ID)
  }).
  On(UserUpdated,
   func(ctx context.Context, data UserEvent) {
    // handle user updated event
    logger.Printf("user updated: %s", data.ID)
   },
   func(ctx context.Context, data UserEvent) {
    // other event listener for the same event.
    logger.Printf("making some actions with user: %s", data.ID)
   },
  )

 emitter.Emit(ctx, UserCreated, UserEvent{ID: "1"})
 emitter.Emit(ctx, UserDeleted, UserEvent{ID: "2"})
 emitter.Emit(ctx, UserUpdated, UserEvent{ID: "3"})

 time.Sleep(1 * time.Second)
}
```

## Bootstrap complex event listeners

In some scenarios, you may need to bootstrap complex event listeners. For example, the listerner
has several dependencies, like a logger, a database connection, etc. In this case, you can use a
function to bootstrap the event listener and pass the dependencies to it.

```go
package main

import (
 "context"
 "log"
 "os"
 "time"

 "github.com/jferrl/go-event"
)

// UserCreated .
const UserCreated event.Event = "user.created"

// UserEvent .
type UserEvent struct {
 ID string
}

// EmailClient handles the email scheduling.
type EmailClient struct{}

// Schedule schedules an email to be sent to the user.
func (e *EmailClient) Schedule(_ context.Context, _ string) error {
 return nil
}

// EmailScheduler handles the email scheduling.
type EmailScheduler interface {
 Schedule(ctx context.Context, userID string) error
}

// Logger handles the logging.
type Logger interface {
 Printf(format string, v ...interface{})
}

// BootstrapEmailSheduleHandler bootstraps the email scheduler handler
// Handler will be called when user created event is emitted.
// It will schedule an email to be sent to the user.
func BootstrapEmailSheduleHandler(logger Logger, s EmailScheduler) event.Listerner[UserEvent] {
 return func(ctx context.Context, data UserEvent) {
  // handle user created event
  logger.Printf("user created: %s", data.ID)

  // schedule email
  if err := s.Schedule(ctx, data.ID); err != nil {
   logger.Printf("failed to schedule email: %v", err)
  }
 }
}

func main() {
 logger := log.New(os.Stdout, "go-event: ", log.LstdFlags)
 e := &EmailClient{}

 ctx := context.Background()

 emitter := event.NewEmitter[UserEvent]()

 emitter.On(UserCreated, BootstrapEmailSheduleHandler(logger, e))

 emitter.Emit(ctx, UserCreated, UserEvent{ID: "1"})

 time.Sleep(1 * time.Second)
}
```

## License

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
