package frontend

import "github.com/RacoonMediaServer/rms-library/pkg/selector"

type Setup struct {
	DiscoveryService DiscoveryService
	TorrentService   TorrentService
	Selector         selector.MediaSelector
	SelectCriterion  selector.Criteria
}
