package aggregation

type Counters interface {
	getLastOffset() int64

	// updateItemsCount will update the counts and return the result with
	// total values for input counts.
	updateItemsCount(lastOffset int64, increments map[string]int64)
}

type counters struct {
	lastOffset   int64
	itemCounters map[string]int64
}

func (c *counters) getLastOffset() int64 {
	return 0
}

func (c *counters) updateItemsCount(lastOffset int64, increments map[string]int64) {
	// TODO: We may have to potentially synchronize
	c.lastOffset = lastOffset
	for itemID, increment := range increments {
		existingVal := c.itemCounters[itemID]
		c.itemCounters[itemID] = existingVal + increment
	}
}

type CountersFactory interface {
	NewCounters() Counters
}

type CountersFactoryFunc func() Counters

func (c CountersFactoryFunc) NewCounters() Counters {
	return c()
}

var _ CountersFactory = CountersFactoryFunc(nil)

func NewCounters() Counters {
	return &counters{
		itemCounters: make(map[string]int64),
	}
}
