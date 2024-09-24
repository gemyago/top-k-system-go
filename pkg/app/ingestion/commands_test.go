package ingestion

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommands(t *testing.T) {
	newMockDeps := func(t *testing.T) CommandsDeps {
		return CommandsDeps{
			ItemEventsWriter: services.NewMockKafkaWriter(t),
		}
	}

	t.Run("IngestItemEvent", func(t *testing.T) {
		t.Run("should write event to the topic", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)
			wantEvt := models.MakeRandomItemEvent()

			body := lo.Must(json.Marshal(&wantEvt))

			mockWriter, _ := mockDeps.ItemEventsWriter.(*services.MockKafkaWriter)
			mockWriter.EXPECT().WriteMessages(
				mock.AnythingOfType("withoutCancelCtx"),
				kafka.Message{
					Key:   []byte(wantEvt.ItemID),
					Value: body,
				},
			).Return(nil)

			lo.Must0(commands.IngestItemEvent(context.Background(), &wantEvt))
		})
		t.Run("should fail if write fails", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			commands := NewCommands(mockDeps)
			wantEvt := models.MakeRandomItemEvent()

			mockWriter, _ := mockDeps.ItemEventsWriter.(*services.MockKafkaWriter)
			wantErr := errors.New(faker.Sentence())
			mockWriter.EXPECT().WriteMessages(
				mock.Anything,
				mock.Anything,
			).Return(wantErr)

			gotErr := commands.IngestItemEvent(context.Background(), &wantEvt)
			require.ErrorIs(t, gotErr, wantErr)
		})
	})
}
