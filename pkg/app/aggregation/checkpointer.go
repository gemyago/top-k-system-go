package aggregation

import (
	"context"
	"fmt"

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

func (cp *checkPointer) dumpState(ctx context.Context, counters Counters) error {
	countersFileName := fmt.Sprintf("counters-%d", counters.getLastOffset())
	newManifest := checkPointManifest{
		LastOffset:           counters.getLastOffset(),
		CountersBlobFileName: countersFileName,
	}
	if err := cp.CheckPointerModel.writeCounters(ctx, countersFileName, counters.getItemsCounters()); err != nil {
		return err
	}
	if err := cp.CheckPointerModel.writeManifest(ctx, newManifest); err != nil {
		return err
	}
	return nil
}

func NewCheckPointer(deps CheckPointerDeps) CheckPointer {
	return &checkPointer{
		CheckPointerDeps: deps,
	}
}
