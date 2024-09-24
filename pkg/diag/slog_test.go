package diag

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLogAttributesFromContext(t *testing.T) {
	t.Run("return empty value if no attributes", func(t *testing.T) {
		got := GetLogAttributesFromContext(context.Background())
		assert.Equal(t, LogAttributes{}, got)
	})
	t.Run("return actual value", func(t *testing.T) {
		want := LogAttributes{CorrelationID: slog.StringValue(faker.UUIDHyphenated())}
		ctx := context.WithValue(context.Background(), contextDiagAttrs, want)
		got := GetLogAttributesFromContext(ctx)
		assert.Equal(t, want, got)
	})
}

func TestSetLogAttributesToContext(t *testing.T) {
	want := LogAttributes{CorrelationID: slog.StringValue(faker.UUIDHyphenated())}
	ctx := SetLogAttributesToContext(context.Background(), want)
	got := GetLogAttributesFromContext(ctx)
	assert.Equal(t, want, got)
}

func TestDiagSlogHandler(t *testing.T) {
	t.Run("WithAttrs", func(t *testing.T) {
		t.Run("should delegate to target", func(t *testing.T) {
			target := NewMockSlogHandler(t)
			mockResult := NewMockSlogHandler(t)
			handler := diagLogHandler{target: target}
			attrs := []slog.Attr{slog.String(faker.Word(), faker.Word())}

			target.EXPECT().WithAttrs(attrs).Return(mockResult)
			assert.Equal(t, mockResult, handler.WithAttrs(attrs))
		})
	})
	t.Run("Handle", func(t *testing.T) {
		t.Run("should delegate to target", func(t *testing.T) {
			target := NewMockSlogHandler(t)
			handler := diagLogHandler{target: target}
			ctx := context.Background()
			originalRec := slog.NewRecord(time.Now(), slog.LevelInfo, faker.Sentence(), 0)
			target.EXPECT().Handle(ctx, originalRec).Return(nil)
			assert.NoError(t, handler.Handle(ctx, originalRec))
		})
		t.Run("should add diag attributes", func(t *testing.T) {
			target := NewMockSlogHandler(t)

			handler := diagLogHandler{target: target}
			attrs := LogAttributes{
				CorrelationID: slog.StringValue(faker.UUIDHyphenated()),
			}
			originalRec := slog.NewRecord(time.Now(), slog.LevelInfo, faker.Sentence(), 0)
			ctx := SetLogAttributesToContext(context.Background(), attrs)
			wantRec := originalRec.Clone()
			wantRec.AddAttrs(slog.Attr{Key: "correlationId", Value: attrs.CorrelationID})
			target.EXPECT().Handle(ctx, wantRec).Return(nil)
			assert.NoError(t, handler.Handle(ctx, originalRec))
		})
	})
	t.Run("SetupRootLogger", func(t *testing.T) {
		t.Run("should setup text handler by default", func(t *testing.T) {
			logger := SetupRootLogger(NewRootLoggerOpts())
			diagHandler, ok := logger.Handler().(*diagLogHandler)
			require.True(t, ok)
			assert.IsType(t, &slog.TextHandler{}, diagHandler.target)
		})
		t.Run("should optionally setup json handler", func(t *testing.T) {
			logger := SetupRootLogger(NewRootLoggerOpts().WithJSONLogs(true).WithLogLevel(slog.LevelDebug))
			diagHandler, ok := logger.Handler().(*diagLogHandler)
			require.True(t, ok)
			assert.IsType(t, &slog.JSONHandler{}, diagHandler.target)
		})
	})
}
