package main

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/app/ingestion"
	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func Test_eventsSender(t *testing.T) {
	newSender := func(t *testing.T) *eventsSenderImpl {
		return &eventsSenderImpl{
			RootLogger:        diag.RootTestLogger(),
			IngestionCommands: ingestion.NewMockCommands(t),
		}
	}

	t.Run("sendTestEvent", func(t *testing.T) {
		t.Run("should send a given number of test events", func(t *testing.T) {
			sender := newSender(t)
			ctx := context.Background()
			wantTestEvents := rand.IntN(5) + 5
			wantItemID := faker.UUIDHyphenated()

			commands, _ := sender.IngestionCommands.(*ingestion.MockCommands)
			wantEvt := &models.ItemEvent{ItemID: wantItemID}
			for range wantTestEvents {
				commands.EXPECT().IngestItemEvent(ctx, wantEvt).Return(nil)
			}

			require.NoError(t, sender.sendTestEvent(ctx, wantItemID, wantTestEvents))
		})
	})
}
