package aggregation

import (
	"container/heap"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/google/btree"
)

// topKGetAllItemsLimit is used to get all items in the topKItems.
const topKGetAllItemsLimit = -1

// topKMaxItemsSize is the maximum number of items that can be stored in the
// topKItems.
const topKMaxItemsSize = 1000

type topKItem struct {
	ItemID string
	Count  int64
}

// String is used mostly for debugging purposes.
func (i *topKItem) String() string {
	return i.ItemID + ":" + strconv.FormatInt(i.Count, 10)
}

type topKItems interface {
	load(vals []*topKItem)
	getItems(limit int) []*topKItem
	updateIfGreater(item topKItem)
}

type topKBTreeItems struct {
	maxSize   int
	tree      *btree.BTreeG[*topKItem]
	itemsByID map[string]*topKItem
}

// getItems returns all items in the tree in descending order.
func (items *topKBTreeItems) getItems(limit int) []*topKItem {
	if limit == topKGetAllItemsLimit {
		limit = items.tree.Len()
	}
	result := make([]*topKItem, 0, limit)
	items.tree.Descend(func(i *topKItem) bool {
		result = append(result, i)
		return len(result) < limit
	})
	return result
}

func (items *topKBTreeItems) load(vals []*topKItem) {
	for _, val := range vals {
		items.tree.ReplaceOrInsert(val)
		items.itemsByID[val.ItemID] = val
	}
	for items.tree.Len() > items.maxSize {
		if minItem, ok := items.tree.Min(); ok {
			items.tree.Delete(minItem)
			delete(items.itemsByID, minItem.ItemID)
		}
	}
}

func (items *topKBTreeItems) updateIfGreater(item topKItem) {
	existingItem := items.itemsByID[item.ItemID]

	// if existing item then we do update only
	if existingItem != nil {
		items.tree.Delete(existingItem)
		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.ItemID] = &item
		return
	}

	// if we have space then we insert
	if items.tree.Len() < items.maxSize {
		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.ItemID] = &item
		return
	}

	// if we don't have space, we only insert if new item is greater than
	// existing min item
	if minItem, ok := items.tree.Min(); ok && item.Count > minItem.Count {
		items.tree.Delete(minItem)
		delete(items.itemsByID, minItem.ItemID)

		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.ItemID] = &item
	}
}

var _ topKItems = (*topKBTreeItems)(nil)

const topKItemsTreeDegree = 10

func newTopKBTreeItems(maxSize int) *topKBTreeItems {
	return &topKBTreeItems{
		maxSize:   maxSize,
		itemsByID: make(map[string]*topKItem),
		tree: btree.NewG(topKItemsTreeDegree, func(a, b *topKItem) bool {
			if a.Count == b.Count {
				return a.ItemID < b.ItemID
			}
			return a.Count < b.Count
		}),
	}
}

type topKHeapItemsList []*topKItem

func (items topKHeapItemsList) Len() int {
	return len(items)
}

func (items topKHeapItemsList) Less(i, j int) bool {
	left := items[i]
	right := items[j]
	return left.Count < right.Count
}

func (items topKHeapItemsList) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items *topKHeapItemsList) Push(x any) {
	*items = append(*items, x.(*topKItem))
}

func (items *topKHeapItemsList) Pop() any {
	n := len(*items)
	x := (*items)[n-1]
	*items = (*items)[:n-1]
	return x
}

var _ heap.Interface = (*topKHeapItemsList)(nil)

type topKHeapItems struct {
	maxSize int
	items   topKHeapItemsList
}

func (items *topKHeapItems) findItemIndex(itemID string) (int, bool) {
	for i, item := range items.items {
		if item.ItemID == itemID {
			return i, true
		}
	}
	return 0, false
}

func (items *topKHeapItems) load(values []*topKItem) {
	items.items = make([]*topKItem, len(values))
	copy(items.items, values)
	heap.Init(&items.items)
	for items.items.Len() > items.maxSize {
		heap.Pop(&items.items)
	}
}

func (items *topKHeapItems) getItems(limit int) []*topKItem {
	if limit == topKGetAllItemsLimit {
		limit = len(items.items)
	}
	result := make([]*topKItem, len(items.items))
	copy(result, items.items)
	slices.SortFunc(result, func(i, j *topKItem) int {
		if i.Count == j.Count {
			return strings.Compare(j.ItemID, i.ItemID)
		}
		return int(j.Count - i.Count)
	})
	if limit >= len(result) {
		return result
	}
	return result[:limit]
}

func (items *topKHeapItems) updateIfGreater(item topKItem) {
	if len(items.items) < items.maxSize {
		heap.Push(&items.items, &item)
		return
	}

	if itemIndex, ok := items.findItemIndex(item.ItemID); ok {
		items.items[itemIndex] = &item
		heap.Fix(&items.items, itemIndex)
		return
	}

	if item.Count > items.items[0].Count {
		heap.Pop(&items.items)
		heap.Push(&items.items, &item)
	}
}

var _ topKItems = (*topKHeapItems)(nil)

func newTopKHeapItems(maxSize int) *topKHeapItems {
	return &topKHeapItems{
		maxSize: maxSize,
		items:   make([]*topKItem, 0, maxSize),
	}
}

type synchronisedTopKItems struct {
	topKItems
	rwLock sync.RWMutex
}

func (items *synchronisedTopKItems) getItems(limit int) []*topKItem {
	items.rwLock.RLock()
	defer items.rwLock.RUnlock()
	return items.topKItems.getItems(limit)
}

func (items *synchronisedTopKItems) updateIfGreater(item topKItem) {
	items.rwLock.Lock()
	defer items.rwLock.Unlock()
	items.topKItems.updateIfGreater(item)
}

func (items *synchronisedTopKItems) load(vals []*topKItem) {
	items.rwLock.Lock()
	defer items.rwLock.Unlock()
	items.topKItems.load(vals)
}

var _ topKItems = (*synchronisedTopKItems)(nil)

type topKItemsFactory interface {
	newTopKItems(maxSize int) topKItems
}

type topKItemsFactoryFunc func(maxSize int) topKItems

func (f topKItemsFactoryFunc) newTopKItems(maxSize int) topKItems {
	return f(maxSize)
}

var _ topKItemsFactory = topKItemsFactoryFunc(nil)

func newTopKItems(maxSize int) topKItems {
	return &synchronisedTopKItems{
		// btree implementation is more performant
		topKItems: newTopKBTreeItems(maxSize),
	}
}
