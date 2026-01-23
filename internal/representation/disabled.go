package representation

import "github.com/racoon-devel/raccoon-pirate/internal/model"

type disabledImpl struct {
}

// Clean implements Service.
func (d *disabledImpl) Clean() {
}

// Register implements Service.
func (d *disabledImpl) Register(t *model.Torrent, location string) {
}

// Unregister implements Service.
func (d *disabledImpl) Unregister(t *model.Torrent) {
}
