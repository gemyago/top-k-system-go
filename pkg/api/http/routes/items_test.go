package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gemyago/top-k-system-go/pkg/app/ingestion"
	"github.com/gemyago/top-k-system-go/pkg/app/models"
	"github.com/gemyago/top-k-system-go/pkg/diag"
	"github.com/gemyago/top-k-system-go/pkg/services"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			Time:       services.NewMockNow(),
		}
		return mockDeps{
			ItemsRoutesDeps: deps,
			Mux:             mux,
		}
	}

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
