package blobstorage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	"go.uber.org/dig"
)

type localStorage struct {
	LocalStorageDeps
	logger *slog.Logger
}

func (s *localStorage) Upload(ctx context.Context, key string, contents io.Reader) error {
	filePath := path.Join(s.LocalStorageFolder, key)
	s.logger.DebugContext(ctx, "Writing file", slog.String("key", key), slog.String("path", filePath))
	file, err := os.Create(filePath)
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

func (s *localStorage) Download(ctx context.Context, key string, out io.Writer) error {
	filePath := path.Join(s.LocalStorageFolder, key)
	s.logger.DebugContext(ctx, "Reading file", slog.String("key", key), slog.String("path", filePath))
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

func (s *localStorage) Delete(_ context.Context, key string) error {
	filePath := path.Join(s.LocalStorageFolder, key)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to remove file %s: %w", filePath, err)
	}
	return nil
}

type LocalStorageDeps struct {
	dig.In

	RootLogger *slog.Logger

	// config
	LocalStorageFolder string `name:"config.blobstorage.localFolder"`
}

func NewLocalStorage(deps LocalStorageDeps) Storage {
	return &localStorage{
		LocalStorageDeps: deps,
		logger:           deps.RootLogger.WithGroup("local-storage"),
	}
}
