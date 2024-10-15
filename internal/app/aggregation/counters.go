package aggregation

type counters interface {
	getItemsCounters() map[string]int64

	getLastOffset() int64

	// updateItemsCount will update the counts and return the result with
	// total values for input counts
	updateItemsCount(lastOffset int64, increments map[string]int64) map[string]int64
}

type countersImpl struct {
	lastOffset   int64
	itemCounters map[string]int64
}

func (c *countersImpl) getItemsCounters() map[string]int64 {
	return c.itemCounters
}

func (c *countersImpl) getLastOffset() int64 {
	return c.lastOffset
}

func (c *countersImpl) updateItemsCount(lastOffset int64, increments map[string]int64) map[string]int64 {
	// TODO: We may have to potentially synchronize
	c.lastOffset = lastOffset
	result := make(map[string]int64, len(increments))
	for itemID, increment := range increments {
		existingVal := c.itemCounters[itemID]
		nextVal := existingVal + increment
		c.itemCounters[itemID] = nextVal
		result[itemID] = nextVal
	}
	return result
}

type countersFactory interface {
	newCounters() counters
}

type countersFactoryFunc func() counters

func (c countersFactoryFunc) newCounters() counters {
	return c()
}

var _ countersFactory = countersFactoryFunc(nil)

func newCounters() counters {
	return &countersImpl{
		itemCounters: make(map[string]int64),
	}
}
