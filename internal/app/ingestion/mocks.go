//go:build !release

package ingestion

import (
	"context"

	"github.com/gemyago/top-k-system-go/internal/app/models"
)

// Mock interfaces are used to generate mock implementations of all of the components
// that will be reused elsewhere in a system. This helps to minimize the amount of
// duplicate mock implementations that need to be written.

type mockCommands interface {
	IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error
}

var _ mockCommands = (*Commands)(nil)
