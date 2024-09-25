package aggregation

import (
	"context"

	"github.com/gemyago/top-k-system-go/pkg/services/blobstorage"
	"go.uber.org/dig"
)

type CheckPointer interface {
	restoreState(ctx context.Context, counters Counters) error
	dumpState(ctx context.Context, counters Counters) error
}

type CheckPointerDeps struct {
	dig.In

	// services
	blobstorage.Storage
}

type checkPointer struct{}

func (cp *checkPointer) restoreState(ctx context.Context, counters Counters) error {
	return nil
}

func (cp *checkPointer) dumpState(ctx context.Context, counters Counters) error {
	return nil
}

func NewCheckPointer(deps CheckPointerDeps) CheckPointer {
	return &checkPointer{}
}
