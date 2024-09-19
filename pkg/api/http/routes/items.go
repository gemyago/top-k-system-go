package routes

import (
	"log/slog"

	"go.uber.org/dig"
)

type ItemsRoutesDeps struct {
	dig.In

	RootLogger *slog.Logger
}

func NewItemsRoutesGroup(_ ItemsRoutesDeps) Group {
	return Group{
		Mount: func(_ router) {

		},
	}
}
