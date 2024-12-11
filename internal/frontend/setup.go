package frontend

import "github.com/racoon-devel/raccoon-pirate/internal/selector"

type Setup struct {
	DiscoveryService DiscoveryService
	TorrentService   TorrentService
	Selector         selector.MediaSelector
	SelectCriterion  selector.Criteria
}
