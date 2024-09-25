package aggregation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gemyago/top-k-system-go/pkg/services/blobstorage"
	"go.uber.org/dig"
)

type checkPointManifest struct {
	LastRevision         int64  `json:"lastRevision"`
	CountersBlobFileName string `json:"countersBlobFileName"`
}

type CheckPointerModel interface {
	readManifest(ctx context.Context) (checkPointManifest, error)
	writeManifest(ctx context.Context, manifest checkPointManifest) error
	readCounters(ctx context.Context, blobFileName string) (map[string]int64, error)
	writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error
}

type CheckPointerModelDeps struct {
	dig.In

	// services
	blobstorage.Storage
}

type checkPointerModel struct {
	CheckPointerModelDeps
}

func (m checkPointerModel) readManifest(ctx context.Context) (checkPointManifest, error) {
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

func (m checkPointerModel) writeManifest(ctx context.Context, manifest checkPointManifest) error {
	var manifestBytes bytes.Buffer
	if err := json.NewEncoder(&manifestBytes).Encode(manifest); err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}
	return m.Storage.Upload(ctx, "manifest.json", &manifestBytes)
}

func (m checkPointerModel) readCounters(ctx context.Context, blobFileName string) (map[string]int64, error) {
	panic("not implemented")
}

func (m checkPointerModel) writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error {
	panic("not implemented")
}

func NewCheckPointerModel(deps CheckPointerModelDeps) CheckPointerModel {
	return &checkPointerModel{CheckPointerModelDeps: deps}
}
