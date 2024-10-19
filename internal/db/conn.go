package db

import (
	"fmt"

	"github.com/racoon-devel/raccoon-pirate/internal/config"
)

func Open(cfg config.Storage) (Database, error) {
	switch cfg.Driver {
	case "cloverdb":
		return newCloverDB(cfg)
	case "json":
		return newJsonDB(cfg)
	}
	return nil, fmt.Errorf("Database driver '%s' is not implemented", cfg.Driver)
}
