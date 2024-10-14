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

type topKItems struct {
	maxSize   int
	tree      *btree.BTreeG[*topKItem]
	itemsByID map[string]*topKItem
}

// getItems returns all items in the tree in descending order.
// TODO: the limit parameter is not used.
func (items *topKItems) getItems(limit int) []*topKItem {
	result := make([]*topKItem, 0, items.tree.Len())
	items.tree.Descend(func(i *topKItem) bool {
		result = append(result, i)
		return true
	})
	return result
}

func (items *topKItems) load(vals []*topKItem) {
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

func (items *topKItems) updateIfGreater(item topKItem) {
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

const topKItemsTreeDegree = 10 // TODO: needs benchmark

func newTopKItems(maxSize int) *topKItems {
	return &topKItems{
		maxSize:   maxSize,
		itemsByID: make(map[string]*topKItem),
		tree: btree.NewG(topKItemsTreeDegree, func(a, b *topKItem) bool {
			return a.count < b.count
		}),
	}
}
