package main

import (
	"bytes"
	"context"
	"io"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/app/ingestion"
	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/gemyago/top-k-system-go/internal/services/blobstorage"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_eventsSender(t *testing.T) {
	newSender := func(t *testing.T) *eventsSenderImpl {
		return &eventsSenderImpl{
			RootLogger:        diag.RootTestLogger(),
			IngestionCommands: ingestion.NewMockCommands(t),
			Storage:           blobstorage.NewMockStorage(t),
			Time:              services.NewMockNow(),
			RandIntN: func(n int) int {
				return n
			},
		}
	}

	t.Run("sendTestEvent", func(t *testing.T) {
		t.Run("should send a given number of test events", func(t *testing.T) {
			sender := newSender(t)
			ctx := context.Background()
			wantTestEvents := rand.IntN(5) + 5
			wantItemID := faker.UUIDHyphenated()

			mockTime, _ := sender.Time.(*services.MockNow)

			commands, _ := sender.IngestionCommands.(*ingestion.MockCommands)
			wantEvt := &models.ItemEvent{
				ItemID:     wantItemID,
				IngestedAt: mockTime.Now(),
			}
			for range wantTestEvents {
				commands.EXPECT().IngestItemEvent(ctx, wantEvt).Return(nil)
			}

			require.NoError(t, sender.sendTestEvent(ctx, wantItemID, wantTestEvents))
		})
	})

	t.Run("sendTestEvents", func(t *testing.T) {
		t.Run("should read events from file and send test events from each item", func(t *testing.T) {
			sender := newSender(t)
			ctx := context.Background()
			wantTestEventsMin := rand.IntN(5)
			wantTestEventsMax := rand.IntN(5) + wantTestEventsMin
			wantItemsFile := faker.DomainName()
			itemIDs := []string{faker.UUIDHyphenated(), faker.UUIDHyphenated(), faker.UUIDHyphenated()}
			wantRandTimes := rand.IntN(5) + 5

			sender.RandIntN = func(n int) int {
				assert.Equal(t, wantTestEventsMax-wantTestEventsMin, n)
				return wantRandTimes
			}

			mockTime, _ := sender.Time.(*services.MockNow)

			storage, _ := sender.Storage.(*blobstorage.MockStorage)
			storage.EXPECT().Download(ctx, wantItemsFile, mock.Anything).RunAndReturn(
				func(_ context.Context, _ string, w io.Writer) error {
					_, err := io.Copy(w, bytes.NewBufferString(strings.Join(itemIDs, "\n")))
					return err
				},
			)

			commands, _ := sender.IngestionCommands.(*ingestion.MockCommands)
			for _, itemID := range itemIDs {
				wantEvt := &models.ItemEvent{ItemID: itemID, IngestedAt: mockTime.Now()}
				for range wantRandTimes {
					commands.EXPECT().IngestItemEvent(ctx, wantEvt).Return(nil)
				}
			}

			require.NoError(t, sender.sendTestEvents(ctx, wantItemsFile, wantTestEventsMin, wantTestEventsMax))
		})
	})
}
