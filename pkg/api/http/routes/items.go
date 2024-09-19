package routes

import (
	"log/slog"
	"net/http"

	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"go.uber.org/dig"
)

type ItemsRoutesDeps struct {
	dig.In

	RootLogger *slog.Logger

	// app layer
	Commands ingestion.Commands

	// service layer
	Time services.TimeProvider
}

func NewItemsRoutesGroup(deps ItemsRoutesDeps) Group {
	commands := deps.Commands
	logger := deps.RootLogger.WithGroup("items-routes")
	return Group{
		Mount: func(r router) {
			r.Handle("POST /items/events/{itemID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				itemID := r.PathValue("itemID")
				err := commands.IngestItemEvent(r.Context(), &models.ItemEvent{
					ItemID:     itemID,
					IngestedAt: deps.Time.Now(),
				})
				if err != nil {
					logger.ErrorContext(r.Context(), "Failed to ingest item event", slog.String("itemID", itemID))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusAccepted)
			}))
		},
	}
}
