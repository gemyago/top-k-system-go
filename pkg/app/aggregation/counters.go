package aggregation

type Counters interface {
	// UpdateItemsCount will update the counts and return the result with
	// total values for input counts. Concurrency safe.
	UpdateItemsCount(newCounts map[string]int) map[string]int
}
