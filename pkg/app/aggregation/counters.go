package aggregation

type Counters interface {
	// UpdateItemsCount will update the counts and return the result with
	// total values for input counts. Concurrency safe.
	UpdateItemsCount(lastOffset int64, newCounts map[string]int) map[string]int
}

type counters struct {
}

func (c *counters) UpdateItemsCount(_ int64, _ map[string]int) map[string]int {
	panic("not implemented")
}

func NewCounters() Counters {
	return &counters{}
}
