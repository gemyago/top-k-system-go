package aggregation

import (
	"context"

	"go.uber.org/dig"
)

type CheckPointer interface {
	restoreState(ctx context.Context, counters Counters) error
	dumpState(ctx context.Context, counters Counters) error
}

type CheckPointerDeps struct {
	dig.In
}

type checkPointer struct{}

func (cp *checkPointer) restoreState(ctx context.Context, counters Counters) error {
	panic("not implemented")
}

func (cp *checkPointer) dumpState(ctx context.Context, counters Counters) error {
	panic("not implemented")
}

func NewCheckPointer(deps CheckPointerDeps) CheckPointer {
	return &checkPointer{}
}
