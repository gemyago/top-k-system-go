package blobstorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"go.uber.org/dig"
)

type localStorage struct {
	LocalStorageDeps
}

func (s *localStorage) Upload(_ context.Context, key string, contents io.Reader) error {
	file, err := os.Create(path.Join(s.LocalStorageFolder, key))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", key, err)
	}
	defer file.Close()

	_, err = io.Copy(file, contents)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", key, err)
	}
	return nil
}

func (s *localStorage) Download(_ context.Context, key string, out io.Writer) error {
	filePath := path.Join(s.LocalStorageFolder, key)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err = io.Copy(out, file); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}

type LocalStorageDeps struct {
	dig.In

	// config
	LocalStorageFolder string `name:"config.blobstorage.localFolder"`
}

func NewLocalStorage(deps LocalStorageDeps) Storage {
	return &localStorage{LocalStorageDeps: deps}
}
