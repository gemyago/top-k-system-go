package aggregation

import (
	"container/heap"
	"slices"
	"strconv"

	"github.com/google/btree"
)

type topKItem struct {
	itemID string
	count  int64
}

// String is used mostly for debugging purposes.
func (i *topKItem) String() string {
	return i.itemID + ":" + strconv.FormatInt(i.count, 10)
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
		items.itemsByID[val.itemID] = val
	}
	for items.tree.Len() > items.maxSize {
		if minItem, ok := items.tree.Min(); ok {
			items.tree.Delete(minItem)
			delete(items.itemsByID, minItem.itemID)
		}
	}
}

func (items *topKBTreeItems) updateIfGreater(item topKItem) {
	existingItem := items.itemsByID[item.itemID]

	// if existing item then we do update only
	if existingItem != nil {
		items.tree.Delete(existingItem)
		delete(items.itemsByID, existingItem.itemID)

		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.itemID] = &item
		return
	}

	// if we have space then we insert
	if items.tree.Len() < items.maxSize {
		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.itemID] = &item
		return
	}

	// if we don't have space, we only insert if new item is greater than
	// existing min item
	if minItem, ok := items.tree.Min(); ok && item.count > minItem.count {
		items.tree.Delete(minItem)
		delete(items.itemsByID, minItem.itemID)

		items.tree.ReplaceOrInsert(&item)
		items.itemsByID[item.itemID] = &item
	}
}

var _ topKItems = (*topKBTreeItems)(nil)

const topKItemsTreeDegree = 10

func newTopKBTreeItems(maxSize int) *topKBTreeItems {
	return &topKBTreeItems{
		maxSize:   maxSize,
		itemsByID: make(map[string]*topKItem),
		tree: btree.NewG(topKItemsTreeDegree, func(a, b *topKItem) bool {
			return a.count < b.count
		}),
	}
}

type topKHeapItemsList []*topKItem

func (items topKHeapItemsList) Len() int {
	return len(items)
}

func (items topKHeapItemsList) Less(i, j int) bool {
	return items[i].count < items[j].count
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
		if item.itemID == itemID {
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
	result := make([]*topKItem, len(items.items))
	copy(result, items.items)
	slices.SortFunc(result, func(i, j *topKItem) int {
		return int(j.count - i.count)
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

	if itemIndex, ok := items.findItemIndex(item.itemID); ok {
		items.items[itemIndex] = &item
		heap.Fix(&items.items, itemIndex)
		return
	}

	if item.count > items.items[0].count {
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
