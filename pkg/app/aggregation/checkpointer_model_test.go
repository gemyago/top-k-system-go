package aggregation

import (
	"context"
	"encoding/json"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/services/blobstorage"
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

	randomManifest := func() checkPointManifest {
		return checkPointManifest{
			LastRevision:         rand.Int64N(10000),
			CountersBlobFileName: faker.Word(),
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
}
