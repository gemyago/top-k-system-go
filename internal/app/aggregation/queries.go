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
	Items []TopKItem `json:"items"`
}

func (q *Queries) GetTopKItems(
	ctx context.Context,
	params GetTopKItemsParams,
) (*GetTopKItemsResponse, error) {
	return nil, nil
}

type QueriesDeps struct {
	dig.In

	topKItemsFactory topKItemsFactory
}

func NewQueries(deps QueriesDeps) *Queries {
	return &Queries{
		allTimeItems: deps.topKItemsFactory.newTopKItems(topKMaxItemsSize),
	}
}
