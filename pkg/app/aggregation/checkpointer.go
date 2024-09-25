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

	// app layer
	CheckPointerModel
}

type checkPointer struct {
	CheckPointerDeps
}

func (cp *checkPointer) restoreState(ctx context.Context, counters Counters) error {
	manifest, err := cp.CheckPointerModel.readManifest(ctx)
	if err != nil {
		return err
	}
	values, err := cp.CheckPointerModel.readCounters(ctx, manifest.CountersBlobFileName)
	if err != nil {
		return err
	}
	counters.updateItemsCount(manifest.LastOffset, values)
	return nil
}

func (cp *checkPointer) dumpState(_ context.Context, _ Counters) error {
	return nil
}

func NewCheckPointer(deps CheckPointerDeps) CheckPointer {
	return &checkPointer{
		CheckPointerDeps: deps,
	}
}
