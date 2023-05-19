//go:generate mockgen -source=interface.go -destination=mocks/mock_iface.go -package=mocks
package depth

import "stock/updater/models"

type Options struct {
	Descending bool
}

type Depth interface {
	SetOptions(opts Options)
	Store([]models.Item)
	Load() []models.Item
	Update(models.Item)
}
