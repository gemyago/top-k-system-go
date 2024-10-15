package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/top-k-system-go/internal/app/aggregation"
	"github.com/gemyago/top-k-system-go/internal/app/ingestion"
	"github.com/gemyago/top-k-system-go/internal/app/models"
	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/gemyago/top-k-system-go/internal/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestItemsRoutes(t *testing.T) {
	type mockDeps struct {
		ItemsRoutesDeps
		Mux *http.ServeMux
	}
	makeDeps := func(t *testing.T) mockDeps {
		mux := http.NewServeMux()
		deps := ItemsRoutesDeps{
			RootLogger: diag.RootTestLogger(),
			Commands:   ingestion.NewMockCommands(t),
			Queries:    aggregation.NewMockQueries(t),
			Time:       services.NewMockNow(),
		}
		return mockDeps{
			ItemsRoutesDeps: deps,
			Mux:             mux,
		}
	}

	t.Run("GET /items/top", func(t *testing.T) {
		t.Run("should return top items", func(t *testing.T) {
			wantLimit := 100 + rand.IntN(100)
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/items/top?limit=%d", wantLimit),
				http.NoBody,
			)
			w := httptest.NewRecorder()
			deps := makeDeps(t)

			mockQueries, _ := deps.Queries.(*aggregation.MockQueries)

			wantResponse := &aggregation.GetTopKItemsResponse{
				Data: []aggregation.TopKItem{
					{
						ItemID: faker.UUIDHyphenated(),
						Count:  rand.Int64N(100),
					},
					{
						ItemID: faker.UUIDHyphenated(),
						Count:  rand.Int64N(100),
					},
					{
						ItemID: faker.UUIDHyphenated(),
						Count:  rand.Int64N(100),
					},
				},
			}

			mockQueries.EXPECT().GetTopKItems(
				mock.AnythingOfType("backgroundCtx"),
				aggregation.GetTopKItemsParams{
					Limit: wantLimit,
				},
			).Return(wantResponse, nil)

			NewItemsRoutesGroup(deps.ItemsRoutesDeps).Mount(deps.Mux)
			deps.Mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			var gotResponse aggregation.GetTopKItemsResponse
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &gotResponse))
		})

		t.Run("should fail if no limit", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/items/top", http.NoBody)
			w := httptest.NewRecorder()
			deps := makeDeps(t)

			NewItemsRoutesGroup(deps.ItemsRoutesDeps).Mount(deps.Mux)
			deps.Mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})

		t.Run("should handle query error", func(t *testing.T) {
			wantLimit := 100 + rand.IntN(100)
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/items/top?limit=%d", wantLimit),
				http.NoBody,
			)
			w := httptest.NewRecorder()
			deps := makeDeps(t)

			mockQueries, _ := deps.Queries.(*aggregation.MockQueries)

			wantErr := errors.New(faker.Sentence())

			mockQueries.EXPECT().GetTopKItems(
				mock.AnythingOfType("backgroundCtx"),
				aggregation.GetTopKItemsParams{
					Limit: wantLimit,
				},
			).Return(nil, wantErr)

			NewItemsRoutesGroup(deps.ItemsRoutesDeps).Mount(deps.Mux)
			deps.Mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})

	t.Run("POST /items/events", func(t *testing.T) {
		t.Run("should ingest the event", func(t *testing.T) {
			wantItemID := faker.UUIDHyphenated()
			req := httptest.NewRequest(http.MethodPost, "/items/events/"+wantItemID, http.NoBody)
			w := httptest.NewRecorder()
			deps := makeDeps(t)

			mockCommands, _ := deps.Commands.(*ingestion.MockCommands)
			mockTime, _ := deps.Time.(*services.MockNow)

			mockCommands.EXPECT().IngestItemEvent(
				mock.AnythingOfType("backgroundCtx"),
				&models.ItemEvent{ItemID: wantItemID, IngestedAt: mockTime.Now()},
			).Return(nil)

			NewItemsRoutesGroup(deps.ItemsRoutesDeps).Mount(deps.Mux)
			deps.Mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusAccepted, w.Code)
		})
		t.Run("should handle error", func(t *testing.T) {
			wantItemID := faker.UUIDHyphenated()
			req := httptest.NewRequest(http.MethodPost, "/items/events/"+wantItemID, http.NoBody)
			w := httptest.NewRecorder()
			deps := makeDeps(t)

			mockCommands, _ := deps.Commands.(*ingestion.MockCommands)
			mockTime, _ := deps.Time.(*services.MockNow)

			mockCommands.EXPECT().IngestItemEvent(
				mock.AnythingOfType("backgroundCtx"),
				&models.ItemEvent{ItemID: wantItemID, IngestedAt: mockTime.Now()},
			).Return(errors.New(faker.Sentence()))

			NewItemsRoutesGroup(deps.ItemsRoutesDeps).Mount(deps.Mux)
			deps.Mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}
