package blobstorage

import (
	"bytes"
	"context"
	"os"
	"path"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalStorage(t *testing.T) {
	newMockDeps := func(t *testing.T) LocalStorageDeps {
		return LocalStorageDeps{
			LocalStorageFolder: t.TempDir(),
		}
	}

	t.Run("upload", func(t *testing.T) {
		t.Run("should write file to the given folder", func(t *testing.T) {
			deps := newMockDeps(t)
			storage := NewLocalStorage(deps)
			ctx := context.Background()
			wantData := faker.Sentence()
			var contents bytes.Buffer
			lo.Must1(contents.WriteString(wantData))
			key := faker.UUIDHyphenated()
			require.NoError(t, storage.Upload(ctx, key, &contents))

			gotData, err := os.ReadFile(path.Join(deps.LocalStorageFolder, key))
			require.NoError(t, err)
			assert.Equal(t, wantData, string(gotData))
		})
	})

	t.Run("download", func(t *testing.T) {
		t.Run("should read given file", func(t *testing.T) {
			deps := newMockDeps(t)
			storage := NewLocalStorage(deps)
			ctx := context.Background()
			wantData := faker.Sentence()
			key := faker.UUIDHyphenated()
			require.NoError(t, os.WriteFile(path.Join(deps.LocalStorageFolder, key), []byte(wantData), 0644))

			var result bytes.Buffer
			require.NoError(t, storage.Download(ctx, key, &result))
			assert.Equal(t, wantData, string(result.String()))
		})
	})
}
