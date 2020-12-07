package types

import (
	"context"
	"github.com/thejerf/abtime"
)

type contextKeyType int

const contextKeyClock contextKeyType = iota

func ContextWithClock(ctx context.Context, time abtime.AbstractTime) context.Context {
	return context.WithValue(ctx, contextKeyClock, time)
}

func ClockFromContext(ctx context.Context) abtime.AbstractTime {
	iface := ctx.Value(contextKeyClock)
	if iface == nil {
		return nil
	}
	return iface.(abtime.AbstractTime)
}

func NewClock(ctx context.Context) abtime.AbstractTime {
	clock := ClockFromContext(ctx)
	if clock == nil {
		clock = NewRealClock()
	}
	return clock
}

func NewRealClock() abtime.RealTime {
	return abtime.NewRealTime()
}

func NewMockClock() *abtime.ManualTime {
	return abtime.NewManual()
}
