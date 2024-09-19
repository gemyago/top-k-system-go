package ingestion

import (
	"context"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
)

type Commands interface {
	IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error
}
