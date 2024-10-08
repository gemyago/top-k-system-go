package aggregation

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckPointerModel(t *testing.T) {
	newMockDeps := func(t *testing.T) CheckPointerModelDeps {
		return CheckPointerModelDeps{
			Storage: blobstorage.NewMockStorage(t),
		}
	}

	t.Run("readManifest", func(t *testing.T) {
		t.Run("should load the manifest from blob storage", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			ctx := context.Background()

			wantManifest := randomManifest()
			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, "manifest.json", mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, w io.Writer) error {
				return json.NewEncoder(w).Encode(&wantManifest)
			})

			gotManifest, err := model.readManifest(ctx)
			require.NoError(t, err)
			assert.Equal(t, wantManifest, gotManifest)
		})
		t.Run("should return error if failed to read manifest", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			ctx := context.Background()
			wantErr := errors.New(faker.Sentence())

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, "manifest.json", mock.Anything,
			).Return(wantErr)

			_, err := model.readManifest(ctx)
			require.ErrorIs(t, err, wantErr)
		})
		t.Run("should return error if failed to decode manifest", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, "manifest.json", mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, w io.Writer) error {
				_, err := w.Write([]byte(faker.Sentence()))
				return err
			})

			_, err := model.readManifest(ctx)
			require.Error(t, err)
		})
	})

	t.Run("writeManifest", func(t *testing.T) {
		t.Run("should upload manifest to blob storage", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			ctx := context.Background()

			wantManifest := randomManifest()
			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Upload(
				ctx, "manifest.json", mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, r io.Reader) error {
				var got checkPointManifest
				require.NoError(t, json.NewDecoder(r).Decode(&got))
				assert.Equal(t, wantManifest, got)
				return nil
			})

			require.NoError(t, model.writeManifest(ctx, wantManifest))
		})
	})

	t.Run("readCounters", func(t *testing.T) {
		t.Run("should read counters from a given file", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			wantCounters := map[string]int64{
				faker.UUIDHyphenated(): rand.Int64(),
				faker.UUIDHyphenated(): rand.Int64(),
				faker.UUIDHyphenated(): rand.Int64(),
			}
			wantFile := faker.Word()

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, wantFile, mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, w io.Writer) error {
				return json.NewEncoder(w).Encode(wantCounters)
			})

			got, err := model.readCounters(ctx, wantFile)
			require.NoError(t, err)
			assert.Equal(t, wantCounters, got)
		})
		t.Run("should return error if failed to read counters", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			wantFile := faker.Word()
			wantErr := errors.New(faker.Sentence())

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, wantFile, mock.Anything,
			).Return(wantErr)

			_, err := model.readCounters(ctx, wantFile)
			require.ErrorIs(t, err, wantErr)
		})
		t.Run("should return error if failed to decode counters", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			wantFile := faker.Word()

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(
				ctx, wantFile, mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, w io.Writer) error {
				_, err := w.Write([]byte(faker.Sentence()))
				return err
			})

			_, err := model.readCounters(ctx, wantFile)
			require.Error(t, err)
		})
	})

	t.Run("writeCounters", func(t *testing.T) {
		t.Run("should write counters to a given file", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			wantCounters := randomCountersValues()
			wantFile := faker.Word()

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Upload(
				ctx, wantFile, mock.Anything,
			).RunAndReturn(func(_ context.Context, _ string, r io.Reader) error {
				var got map[string]int64
				require.NoError(t, json.NewDecoder(r).Decode(&got))
				assert.Equal(t, wantCounters, got)
				return nil
			})

			err := model.writeCounters(ctx, wantFile, wantCounters)
			require.NoError(t, err)
		})

		t.Run("should return error if failed to upload counters", func(t *testing.T) {
			deps := newMockDeps(t)
			model := NewCheckPointerModel(deps)

			wantCounters := randomCountersValues()
			wantFile := faker.Word()
			wantErr := errors.New(faker.Sentence())

			ctx := context.Background()

			storage, _ := deps.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Upload(
				ctx, wantFile, mock.Anything,
			).Return(wantErr)

			err := model.writeCounters(ctx, wantFile, wantCounters)
			require.ErrorIs(t, err, wantErr)
		})
	})
}
