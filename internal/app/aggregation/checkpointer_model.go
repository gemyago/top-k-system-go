package aggregation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
	"go.uber.org/dig"
)

type checkPointManifest struct {
	LastOffset           int64  `json:"lastOffset"`
	CountersBlobFileName string `json:"countersBlobFileName"`
}

type checkPointerModel interface {
	readManifest(ctx context.Context) (checkPointManifest, error)
	writeManifest(ctx context.Context, manifest checkPointManifest) error
	readCounters(ctx context.Context, blobFileName string) (map[string]int64, error)
	writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error
}

type CheckPointerModelDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	// services
	blobstorage.Storage
}

type checkPointerModelImpl struct {
	CheckPointerModelDeps
}

func (m checkPointerModelImpl) readManifest(ctx context.Context) (checkPointManifest, error) {
	var manifestBytes bytes.Buffer
	if err := m.Storage.Download(ctx, "manifest.json", &manifestBytes); err != nil {
		return checkPointManifest{}, fmt.Errorf("failed to read the manifest: %w", err)
	}
	var manifest checkPointManifest
	if err := json.NewDecoder(&manifestBytes).Decode(&manifest); err != nil {
		return checkPointManifest{}, fmt.Errorf("faield to decode manifest: %w", err)
	}
	return manifest, nil
}

func (m checkPointerModelImpl) writeManifest(ctx context.Context, manifest checkPointManifest) error {
	var manifestBytes bytes.Buffer
	if err := json.NewEncoder(&manifestBytes).Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}
	return m.Storage.Upload(ctx, "manifest.json", &manifestBytes)
}

// TODO: blobs are going to be very large (5GB), we may need to consider chunking approach
// but this is a caller level refactoring very likely
// both read and write sides will need to be updated

func (m checkPointerModelImpl) readCounters(ctx context.Context, blobFileName string) (map[string]int64, error) {
	var contents bytes.Buffer
	//TODO: Use gob instead of json
	if err := m.Storage.Download(ctx, blobFileName, &contents); err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	var result map[string]int64
	if err := json.NewDecoder(&contents).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode counters: %w", err)
	}
	return result, nil
}

func (m checkPointerModelImpl) writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error {
	var contents bytes.Buffer
	if err := json.NewEncoder(&contents).Encode(val); err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}
	if err := m.Storage.Upload(ctx, blobFileName, &contents); err != nil {
		return fmt.Errorf("failed to upload blob file %s: %w", blobFileName, err)
	}
	return nil
}

func newCheckPointerModel(deps CheckPointerModelDeps) checkPointerModel {
	return &checkPointerModelImpl{CheckPointerModelDeps: deps}
}
