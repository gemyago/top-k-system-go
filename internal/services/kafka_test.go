package services

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/gemyago/top-k-system-go/internal/diag"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKafka(t *testing.T) {
	t.Run("ItemEventsKafkaReader", func(t *testing.T) {
		makeMockDeps := func() ItemEventsKafkaReaderDeps {
			return ItemEventsKafkaReaderDeps{
				RootLogger:    diag.RootTestLogger(),
				KafkaAddress:  faker.DomainName(),
				KafkaTopic:    faker.DomainName(),
				ReaderMaxWait: 10 * time.Second,
				ShutdownHooks: NewTestShutdownHooks(),
			}
		}

		t.Run("ReadLastOffset", func(t *testing.T) {
			t.Run("should return the last offset", func(t *testing.T) {
				deps := makeMockDeps()
				mockConn := newMockKafkaConn(t)
				deps.KafkaLeaderDialer = func(
					_ context.Context,
					network, addr, topic string,
					partition int,
				) (kafkaConn, error) {
					assert.Equal(t, "tcp", network)
					assert.Equal(t, deps.KafkaAddress, addr)
					assert.Equal(t, deps.KafkaTopic, topic)
					assert.Equal(t, 0, partition)
					return mockConn, nil
				}
				reader := NewItemEventsKafkaReader(deps)
				ctx := context.Background()
				wantOffset := rand.Int63()

				mockConn.EXPECT().Close().Return(nil)
				mockConn.EXPECT().ReadLastOffset().Return(wantOffset, nil)

				gotOffset, err := reader.ReadLastOffset(ctx)
				require.NoError(t, err)
				assert.Equal(t, wantOffset, gotOffset)
			})

			t.Run("should return error if failed to dial kafka", func(t *testing.T) {
				deps := makeMockDeps()
				wantErr := errors.New(faker.Sentence())
				deps.KafkaLeaderDialer = func(
					_ context.Context,
					_ string, _ string, _ string, _ int,
				) (kafkaConn, error) {
					return nil, wantErr
				}
				reader := NewItemEventsKafkaReader(deps)
				ctx := context.Background()

				_, err := reader.ReadLastOffset(ctx)
				require.Error(t, err)
				assert.ErrorIs(t, err, wantErr)
			})

			t.Run("should return error if failed to read last offset", func(t *testing.T) {
				deps := makeMockDeps()
				mockConn := newMockKafkaConn(t)
				deps.KafkaLeaderDialer = func(
					_ context.Context,
					_ string, _ string, _ string, _ int,
				) (kafkaConn, error) {
					return mockConn, nil
				}
				reader := NewItemEventsKafkaReader(deps)
				ctx := context.Background()
				wantErr := errors.New(faker.Sentence())

				mockConn.EXPECT().Close().Return(nil)
				mockConn.EXPECT().ReadLastOffset().Return(0, wantErr)

				_, err := reader.ReadLastOffset(ctx)
				require.Error(t, err)
				assert.ErrorIs(t, err, wantErr)
			})
		})
	})
}
