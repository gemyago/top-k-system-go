package aggregation

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	"go.uber.org/dig"
)

type aggregationState struct {
	counters     counters
	allTimeItems topKItems
}

type checkPointer interface {
	restoreState(ctx context.Context, state aggregationState) error
	dumpState(ctx context.Context, state aggregationState) error
}

type CheckPointerDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	RootLogger *slog.Logger

	// package private components
	CheckPointerModel checkPointerModel
}

type checkPointerImpl struct {
	logger *slog.Logger
	deps   CheckPointerDeps
}

func (cp *checkPointerImpl) restoreState(ctx context.Context, state aggregationState) error {
	manifest, err := cp.deps.CheckPointerModel.readManifest(ctx)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			cp.logger.InfoContext(ctx, "Manifest not found. No state to restore from.")
			return nil
		}
		return err
	}

	// TODO: read in parallel

	counterValues, err := cp.deps.CheckPointerModel.readCounters(ctx, manifest.CountersBlobFileName)
	if err != nil {
		return fmt.Errorf("failed to read counters: %w", err)
	}
	state.counters.updateItemsCount(manifest.LastOffset, counterValues)

	allTimeItems, err := cp.deps.CheckPointerModel.readItems(ctx, manifest.AllTimeItemsFileName)
	if err != nil {
		return fmt.Errorf("failed to read all time items: %w", err)
	}
	state.allTimeItems.load(allTimeItems)

	return nil
}

func (cp *checkPointerImpl) dumpState(ctx context.Context, state aggregationState) error {
	countersFileName := fmt.Sprintf("counters-%d", state.counters.getLastOffset())
	allTimeItemsFileName := fmt.Sprintf("all-time-items-%d", state.counters.getLastOffset())
	newManifest := checkPointManifest{
		LastOffset:           state.counters.getLastOffset(),
		CountersBlobFileName: countersFileName,
		AllTimeItemsFileName: allTimeItemsFileName,
	}
	// TODO: write in parallel (except the manifest)

	if err := cp.deps.CheckPointerModel.writeCounters(
		ctx,
		countersFileName,
		state.counters.getItemsCounters(),
	); err != nil {
		return fmt.Errorf("failed to write counters: %w", err)
	}

	if err := cp.deps.CheckPointerModel.writeItems(
		ctx,
		allTimeItemsFileName,
		state.allTimeItems.getItems(topKGetAllItemsLimit),
	); err != nil {
		return fmt.Errorf("failed to write all time items: %w", err)
	}

	// We write manifest last so if counters fail, the manifest will point on the last
	// counters
	if err := cp.deps.CheckPointerModel.writeManifest(ctx, newManifest); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}
	return nil
}

func newCheckPointer(deps CheckPointerDeps) checkPointer {
	return &checkPointerImpl{
		logger: deps.RootLogger.WithGroup("check-pointer"),
		deps:   deps,
	}
}
