package types

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/thejerf/abtime"
	"reflect"
	"testing"
)

func TestClockFromContext(t *testing.T) {
	clock := NewClock(context.Background())
	ctx := ContextWithClock(context.Background(), clock)

	tests := []struct {
		name string
		ctx context.Context
		want abtime.AbstractTime
	}{
		{
			name: "ExistsInContext",
			ctx: ctx,
			want: clock,
		},
		{
			name: "NotExistsInContext",
			ctx: context.Background(),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClockFromContext(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClockFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextWithClock(t *testing.T) {
	clock := NewClock(context.Background())
	ctx := ContextWithClock(context.Background(), clock)
	assert.Equal(t, clock, ClockFromContext(ctx))
}

func TestNewClock(t *testing.T) {
	clock := NewClock(context.Background())
	ctx := ContextWithClock(context.Background(), clock)

	tests := []struct {
		name string
		ctx context.Context
		want abtime.AbstractTime
	}{
		{
			name: "ExistsInContext",
			ctx: ctx,
			want: clock,
		},
		{
			name: "NotExistsInContext",
			ctx: context.Background(),
			want: NewRealClock(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClock(tt.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMockClock(t *testing.T) {
	assert.NotNil(t, NewMockClock())
}

func TestNewRealClock(t *testing.T) {
	assert.True(t, reflect.DeepEqual(NewRealClock(), abtime.NewRealTime()))
}
