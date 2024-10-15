package aggregation

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueries(t *testing.T) {
	makeMockDeps := func(t *testing.T) QueriesDeps {
		return QueriesDeps{
			TopKItemsFactory: newMockTopKItemsFactory(t),
		}
	}

	t.Run("GetTopKItems", func(t *testing.T) {
		t.Run("should return all time top k items", func(t *testing.T) {
			deps := makeMockDeps(t)

			mockItems := newMockTopKItems(t)
			mockFactory, _ := deps.TopKItemsFactory.(*mockTopKItemsFactory)
			mockFactory.EXPECT().newTopKItems(topKMaxItemsSize).Return(mockItems)

			wantSize := 10 + rand.IntN(10)
			wantRawItems := randomTopKItems(10)
			mockItems.EXPECT().getItems(wantSize).Return(wantRawItems)

			ctx := context.Background()
			queries := NewQueries(deps)
			got, err := queries.GetTopKItems(ctx, GetTopKItemsParams{Limit: wantSize})
			require.NoError(t, err)

			wantItems := lo.Map(
				wantRawItems,
				func(item *topKItem, _ int) TopKItem {
					return TopKItem{
						ItemID: item.ItemID,
						Count:  item.Count,
					}
				},
			)

			assert.Equal(t, wantItems, got.Data)
		})
	})
}
