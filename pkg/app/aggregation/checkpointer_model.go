package aggregation

import "context"

type checkPointManifest struct {
	LastRevision         int64  `json:"lastRevision"`
	CountersBlobFileName string `json:"countersBlobFileName"`
}

type CheckPointerModel interface {
	readManifest(ctx context.Context) (checkPointManifest, error)
	writeManifest(ctx context.Context, m checkPointManifest) error
	readCounters(ctx context.Context, blobFileName string) (map[string]int64, error)
	writeCounters(ctx context.Context, blobFileName string, val map[string]int64) error
}
