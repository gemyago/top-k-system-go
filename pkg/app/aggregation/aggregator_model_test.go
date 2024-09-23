package aggregation

import (
	"context"
	"encoding/json"
	"io"
	"math/rand"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/samber/lo"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregatorModel(t *testing.T) {
	newMockDeps := func(t *testing.T) ItemEventsAggregatorModelDeps {
		return ItemEventsAggregatorModelDeps{
			ItemEventsReader: services.NewMockKafkaReader(t),
		}
	}

	t.Run("aggregateItemEvent", func(t *testing.T) {
		t.Run("should set new counters to 1", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			model := NewItemEventsAggregatorModel(mockDeps)

			baseOffset := rand.Int63()
			itemEvents := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}
			for i, e := range itemEvents {
				model.aggregateItemEvent(baseOffset+int64(i), &e)
			}

			modelImpl, _ := model.(*itemEventsAggregatorModel)
			assert.Equal(t, baseOffset+int64(len(itemEvents)-1), modelImpl.lastAggregatedOffset)
			for _, e := range itemEvents {
				assert.Equal(t, int64(1), modelImpl.aggregatedItems[e.ItemID])
			}
		})
		t.Run("should increment existing counters", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			model := NewItemEventsAggregatorModel(mockDeps)

			baseCounter := rand.Int63()
			baseOffset := rand.Int63()
			itemEvents := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}
			modelImpl, _ := model.(*itemEventsAggregatorModel)
			for i, e := range itemEvents {
				modelImpl.aggregatedItems[e.ItemID] = baseCounter + int64(i)
				model.aggregateItemEvent(baseOffset+int64(i), &e)
			}

			assert.Equal(t, baseOffset+int64(len(itemEvents)-1), modelImpl.lastAggregatedOffset)
			for i, e := range itemEvents {
				assert.Equal(t, baseCounter+int64(i+1), modelImpl.aggregatedItems[e.ItemID])
			}
		})
	})

	t.Run("fetchMessages", func(t *testing.T) {
		t.Run("should deserialize and feed messages to the channel", func(t *testing.T) {
			mockDeps := newMockDeps(t)
			model := NewItemEventsAggregatorModel(mockDeps)

			baseOffset := rand.Int63()
			itemEvents := []models.ItemEvent{
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
				models.MakeRandomItemEvent(),
			}
			ctx := context.Background()
			mockReader, _ := mockDeps.ItemEventsReader.(*services.MockKafkaReader)

			fetchMessageCounter := 0
			mockReader.EXPECT().FetchMessage(ctx).RunAndReturn(
				func(_ context.Context) (kafka.Message, error) {
					defer func() {
						fetchMessageCounter++
					}()
					if fetchMessageCounter == len(itemEvents) {
						return kafka.Message{}, io.EOF
					}
					nextEvt := itemEvents[fetchMessageCounter]
					data := lo.Must(json.Marshal(nextEvt))
					return kafka.Message{
						Offset: baseOffset + int64(fetchMessageCounter),
						Key:    []byte(nextEvt.ItemID),
						Value:  data,
					}, nil
				},
			)

			gotResults := make([]fetchMessageResult, 0, len(itemEvents))
			syncChan := make(chan struct{})
			go func() {
				for res := range model.fetchMessages(ctx) {
					gotResults = append(gotResults, res)
					if len(gotResults) == len(itemEvents) {
						syncChan <- struct{}{}
						break
					}
				}
			}()
			<-syncChan

			require.Len(t, gotResults, len(itemEvents))
			for i, wantItem := range itemEvents {
				gotResult := gotResults[i]
				assert.Equal(t, fetchMessageResult{
					offset: baseOffset + int64(i),
					event:  &wantItem,
				}, gotResult)
			}
		})
	})
}
