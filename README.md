# go-event

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/jferrl/go-event)
[![Test Status](https://github.com/jferrl/go-event/workflows/tests/badge.svg)](https://github.com/jferrl/go-event/actions?query=workflow%3Atests)
[![codecov](https://codecov.io/gh/jferrl/go-event/branch/main/graph/badge.svg?token=68I4BZF235)](https://codecov.io/gh/jferrl/go-event)
[![Go Report Card](https://goreportcard.com/badge/github.com/jferrl/go-event)](https://goreportcard.com/report/github.com/jferrl/go-event)

Go simple library for event handling.

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

## License

This library is distributed under the MIT license found in the [LICENSE](./LICENSE)
