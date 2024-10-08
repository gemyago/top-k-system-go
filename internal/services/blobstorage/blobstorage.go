package blobstorage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, key string, contents io.Reader) error
	Download(ctx context.Context, key string, out io.Writer) error
	Delete(ctx context.Context, key string) error
}
