package torrents

import (
	"time"

	rms_torrent "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-torrent"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/model"
)

const connectRetryInterval = 10 * time.Second

func (s *Service) trySyncTorrents(torrents []*model.Torrent) {
	for {
		remoteTorrents, err := s.e.List(s.ctx, true)
		if err == nil {
			s.syncTorrents(torrents, remoteTorrents)
			return
		}

		log.Warnf("Cannot connect to torrent server: %s", err)
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(connectRetryInterval):
		}
	}
}

func buildTorrentListIndex(torrents []*rms_torrent.TorrentInfo) map[string]*rms_torrent.TorrentInfo {
	result := map[string]*rms_torrent.TorrentInfo{}
	for _, t := range torrents {
		result[t.Id] = t
	}
	return result
}

func (s *Service) syncTorrents(torrents []*model.Torrent, remoteTorrents []*rms_torrent.TorrentInfo) {
	remoteIdx := buildTorrentListIndex(remoteTorrents)
	for _, t := range torrents {
		remoteInfo, exists := remoteIdx[t.ID]
		if !exists {
			resp, err := s.e.Add(s.ctx, determineCategory(t), t.Title, nil, t.Content)
			if err != nil {
				log.Warnf("Add non existing torrent failed: %s", err)
			} else {
				s.rep.Register(t, resp.Location)
			}
		} else {
			s.rep.Register(t, remoteInfo.Location)
		}
	}

	log.Info("Torrents synced!")
}
