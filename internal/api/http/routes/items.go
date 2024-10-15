package routes

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gemyago/top-k-system-go/internal/app/aggregation"
	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"go.uber.org/dig"
)

type ingestionCommands interface {
	IngestItemEvent(ctx context.Context, evt *models.ItemEvent) error
}

type aggregationQueries interface {
	GetTopKItems(
		_ context.Context,
		params aggregation.GetTopKItemsParams,
	) (*aggregation.GetTopKItemsResponse, error)
}

type ItemsRoutesDeps struct {
	dig.In

	RootLogger *slog.Logger

	// app layer
	Commands ingestionCommands
	Queries  aggregationQueries

	// service layer
	Time services.TimeProvider
}

func NewItemsRoutesGroup(deps ItemsRoutesDeps) Group {
	commands := deps.Commands
	logger := deps.RootLogger.WithGroup("items-routes")
	return Group{
		Mount: func(r router) {
			r.Handle("GET /items/top", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
				if err != nil {
					logger.ErrorContext(r.Context(), "Failed to parse limit", diag.ErrAttr(err))
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				resp, err := deps.Queries.GetTopKItems(r.Context(), aggregation.GetTopKItemsParams{
					Limit: int(limit),
				})
				if err != nil {
					logger.ErrorContext(r.Context(), "Failed to get top items", diag.ErrAttr(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err = json.NewEncoder(w).Encode(resp); err != nil {
					logger.ErrorContext(r.Context(), "Failed to encode response", diag.ErrAttr(err))
					return
				}
			}))
			r.Handle("POST /items/events/{itemID}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				itemID := r.PathValue("itemID")
				err := commands.IngestItemEvent(r.Context(), &models.ItemEvent{
					ItemID:     itemID,
					IngestedAt: deps.Time.Now(),
				})
				if err != nil {
					logger.ErrorContext(r.Context(), "Failed to ingest item event", slog.String("itemID", itemID), diag.ErrAttr(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusAccepted)
			}))
		},
	}
}
