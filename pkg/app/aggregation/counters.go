package aggregation

type Counters interface {
	// updateItemsCount will update the counts and return the result with
	// total values for input counts.
	updateItemsCount(lastOffset int64, increments map[string]int64)
}

type counters struct {
	lastOffset   int64
	itemCounters map[string]int64
}

func (c *counters) updateItemsCount(lastOffset int64, increments map[string]int64) {
	// TODO: We may have to potentially synchronize
	c.lastOffset = lastOffset
	for itemID, increment := range increments {
		existingVal := c.itemCounters[itemID]
		c.itemCounters[itemID] = existingVal + increment
	}
}

func NewCounters() Counters {
	return &counters{
		itemCounters: make(map[string]int64),
	}
}
