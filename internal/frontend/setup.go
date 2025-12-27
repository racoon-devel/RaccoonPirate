package frontend

import "github.com/RacoonMediaServer/rms-library/pkg/selector"

type Setup struct {
	Version                string
	DiscoveryService       DiscoveryService
	SmartSearchService     SmartSearchService
	TorrentService         TorrentService
	Selector               selector.MediaSelector
	SelectCriterion        selector.Criteria
	TelegramAccessProvider TelegramAccessProvider
}
