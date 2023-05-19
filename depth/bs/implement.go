package bs

import (
	"sort"
	"stock/updater/depth"
	"stock/updater/models"
)

// bsStock binary search implementation of the stock
type bsStock struct {
	depth.Options
	items []models.Item
}

func NewDepth(o depth.Options) depth.Depth {
	return &bsStock{
		Options: o,
		items:   make([]models.Item, 1),
	}
}

func (bs *bsStock) SetOptions(opts depth.Options) {
	bs.Options = opts
}

func (bs *bsStock) less(i, j models.Item) bool {
	return (i.Price < j.Price) != bs.Descending
}

func (bs *bsStock) Store(items []models.Item) {
	bs.items = items
	sort.SliceStable(bs.items, func(i, j int) bool {
		return bs.less(bs.items[i], bs.items[j])
	})
}
func (bs *bsStock) Load() []models.Item {
	return bs.items
}

func (bs *bsStock) Update(item models.Item) {
	if len(bs.items) == 0 {
		if item.Size == 0 {
			return
		}
		bs.items = append(bs.items, item)
		return
	}
	pos, toInsert := bs.search(item)
	if item.Size == 0 {
		if toInsert {
			return
		}
		bs.items = append(bs.items[:pos], bs.items[pos+1:]...)
	} else {
		if toInsert {
			if len(bs.items) == pos {
				bs.items = append(bs.items, item)
			} else {
				bs.items = append(bs.items[:pos+1], bs.items[pos:]...)
			}
		}
		bs.items[pos] = item
	}
}

// search finds the position of item in array, otherwise returns the insert position
func (bs *bsStock) search(item models.Item) (pos int, insert bool) {
	l, r := 0, len(bs.items)-1
	for l <= r {
		m := l + (r-l)/2
		val := bs.items[m]
		if val.Price == item.Price {
			return m, false
		}

		if bs.less(val, item) {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return l, true
}
