package event

import (
	"context"
	"testing"
	"time"
)

type eventPayload struct {
	id     string
	called bool
}

func TestEmitter_Emit(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		ev   Event
		data *eventPayload
	}
	tests := []struct {
		name         string
		e            *Emitter[*eventPayload]
		args         args
		listeners    []Listerner[*eventPayload]
		eventEmmited bool
	}{
		{
			name: "empty event listerners",
			e:    NewEmitter[*eventPayload](),
			args: args{
				ctx: ctx,
				ev:  "test",
				data: &eventPayload{
					id: "test",
				},
			},
		},
		{
			name: "single event listerner",
			e:    NewEmitter[*eventPayload](),
			args: args{
				ctx:  ctx,
				ev:   "test",
				data: &eventPayload{},
			},
			listeners: []Listerner[*eventPayload]{
				func(ctx context.Context, data *eventPayload) {
					data.called = true
				},
			},
			eventEmmited: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, l := range tt.listeners {
				tt.e.On(tt.args.ev, l)
			}

			tt.e.Emit(tt.args.ctx, tt.args.ev, tt.args.data)

			// wait for all listeners to finish
			time.Sleep(1 * time.Second)

			if !tt.args.data.called && tt.eventEmmited {
				t.Errorf("data.called = %v, want %v", tt.args.data.called, true)
			}
		})
	}
}
