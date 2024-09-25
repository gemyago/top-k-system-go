package blobstorage

import (
	"context"
	"io"

	"go.uber.org/dig"
)

type localStorage struct {
}

func (s *localStorage) Upload(ctx context.Context, key string, contents io.Reader) error {
	panic("not implemented")
}

func (s *localStorage) Download(ctx context.Context, key string, out io.Writer) error {
	panic("not implemented")
}

type LocalStorageDeps struct {
	dig.In

	// config
	LocalStorageFolder string `name:"config.blobstorage.localFolder"`
}

func NewLocalStorage(deps LocalStorageDeps) Storage {
	return &localStorage{}
}
