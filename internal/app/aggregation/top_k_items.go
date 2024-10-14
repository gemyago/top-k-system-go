package aggregation

import (
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
