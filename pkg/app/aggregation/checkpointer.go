package aggregation

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"go.uber.org/dig"
)

type CheckPointer interface {
	restoreState(ctx context.Context, counters Counters) error
	dumpState(ctx context.Context, counters Counters) error
}

type CheckPointerDeps struct {
	dig.In

	RootLogger *slog.Logger

	// app layer
	CheckPointerModel
}

type checkPointer struct {
	logger *slog.Logger
	CheckPointerDeps
}

func (cp *checkPointer) restoreState(ctx context.Context, counters Counters) error {
	manifest, err := cp.CheckPointerModel.readManifest(ctx)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cp.logger.InfoContext(ctx, "Manifest not found. No state to restore from.")
			return nil
		}
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

	// We write manifest last so if counters fail, the manifest will point on the last
	// counters
	if err := cp.CheckPointerModel.writeManifest(ctx, newManifest); err != nil {
		return err
	}
	return nil
}

func NewCheckPointer(deps CheckPointerDeps) CheckPointer {
	return &checkPointer{
		logger:           deps.RootLogger.WithGroup("check-pointer"),
		CheckPointerDeps: deps,
	}
}
