package services

import "github.com/gofrs/uuid/v5"

type UUIDGenerator func() string

func NewUUIDGenerator() UUIDGenerator {
	return func() string {
		return uuid.Must(uuid.NewV4()).String()
	}
}
