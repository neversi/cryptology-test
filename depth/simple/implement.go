package simple

import (
	"sort"
	"stock/updater/depth"
	"stock/updater/models"
)

type stock struct {
	depth.Options
	items []models.Item
}

func NewDepth() depth.Depth {
	return &stock{}
}

func (s *stock) SetOptions(opts depth.Options) {
	s.Options = opts
}

func (s *stock) Store(items []models.Item) {
	s.items = items
	sort.SliceStable(s.items, func(i, j int) bool {
		return s.less(s.items[i], s.items[j])
	})
}
func (s *stock) Load() []models.Item {
	return s.items
}
func (s *stock) Update(item models.Item) {
	if len(s.items) == 0 {
		s.items = append(s.items, item)
		return
	}
	pos := -1
	for i := range s.items {
		val := s.items[i]
		if val.Price == item.Price {
			pos = i
			break
		}
	}
	switch item.Size {
	case 0:
		if pos == -1 {
			return
		} else {
			s.items = append(s.items[:pos], s.items[pos+1:]...)
		}
	default:
		if pos == -1 {
			s.items = append(s.items, item)
			sort.SliceStable(s.items, func(i, j int) bool {
				return s.less(s.items[i], s.items[j])
			})
		} else {
			s.items[pos].Size = item.Size
		}
	}
}

func (s *stock) less(i, j models.Item) bool {
	return (i.Price < j.Price) != s.Descending
}
