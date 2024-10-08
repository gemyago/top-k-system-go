//go:build !release

package services

import (
	"time"

	"github.com/go-faker/faker/v4"
)

type MockNow struct {
	value time.Time
}

var _ TimeProvider = &MockNow{}

func (m *MockNow) SetValue(t time.Time) {
	m.value = t
}

func (m *MockNow) Now() time.Time {
	return m.value
}

func NewMockNow() *MockNow {
	return &MockNow{
		value: time.UnixMilli(faker.RandomUnixTime()),
	}
}

func MockNowValue(p TimeProvider) time.Time {
	mp, ok := p.(*MockNow)
	if !ok {
		panic("provided TimeProvider is not a MockNow")
	}
	return mp.value
}
