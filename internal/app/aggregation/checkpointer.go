package aggregation

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"go.uber.org/dig"
)

type checkPointer interface {
	restoreState(ctx context.Context, counters counters) error
	dumpState(ctx context.Context, counters counters) error
}

type CheckPointerDeps struct {
	dig.In

	RootLogger *slog.Logger

	// package private components
	checkPointerModel
}

type checkPointerImpl struct {
	logger *slog.Logger
	CheckPointerDeps
}

func (cp *checkPointerImpl) restoreState(ctx context.Context, counters counters) error {
	manifest, err := cp.checkPointerModel.readManifest(ctx)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cp.logger.InfoContext(ctx, "Manifest not found. No state to restore from.")
			return nil
		}
		return err
	}
	values, err := cp.checkPointerModel.readCounters(ctx, manifest.CountersBlobFileName)
	if err != nil {
		return err
	}
	counters.updateItemsCount(manifest.LastOffset, values)
	return nil
}

func (cp *checkPointerImpl) dumpState(ctx context.Context, counters counters) error {
	countersFileName := fmt.Sprintf("counters-%d", counters.getLastOffset())
	newManifest := checkPointManifest{
		LastOffset:           counters.getLastOffset(),
		CountersBlobFileName: countersFileName,
	}
	if err := cp.checkPointerModel.writeCounters(ctx, countersFileName, counters.getItemsCounters()); err != nil {
		return err
	}

	// We write manifest last so if counters fail, the manifest will point on the last
	// counters
	if err := cp.checkPointerModel.writeManifest(ctx, newManifest); err != nil {
		return err
	}
	return nil
}

func newCheckPointer(deps CheckPointerDeps) checkPointer {
	return &checkPointerImpl{
		logger:           deps.RootLogger.WithGroup("check-pointer"),
		CheckPointerDeps: deps,
	}
}
