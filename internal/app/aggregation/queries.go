package aggregation

import (
	"context"

	"go.uber.org/dig"
)

type Queries struct {
	allTimeItems topKItems
}

type GetTopKItemsParams struct {
	Limit int
}

type TopKItem struct {
	ItemID string `json:"itemId"`
	Count  int64  `json:"count"`
}

type GetTopKItemsResponse struct {
	Data []TopKItem `json:"data"`
}

func (q *Queries) GetTopKItems(
	_ context.Context,
	params GetTopKItemsParams,
) (*GetTopKItemsResponse, error) {
	items := q.allTimeItems.getItems(params.Limit)
	result := make([]TopKItem, len(items))
	for i, item := range items {
		result[i] = TopKItem{
			ItemID: item.ItemID,
			Count:  item.Count,
		}
	}
	return &GetTopKItemsResponse{Data: result}, nil
}

type QueriesDeps struct {
	// all injectable fields must be exported
	// to let dig inject them

	dig.In

	TopKItemsFactory topKItemsFactory
}

func NewQueries(deps QueriesDeps) *Queries {
	return &Queries{
		allTimeItems: deps.TopKItemsFactory.newTopKItems(topKMaxItemsSize),
	}
}
