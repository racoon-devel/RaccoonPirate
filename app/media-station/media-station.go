package main

import (
	"flag"
	"github.com/apex/log"
	"github.com/racoon-devel/media-station/internal/config"
	"github.com/racoon-devel/media-station/internal/discovery"
	"github.com/racoon-devel/media-station/internal/selector"
	"github.com/racoon-devel/media-station/internal/torrents"
	"github.com/racoon-devel/media-station/internal/web"
)

var Version = "0.0.0"

func getVoicePriorityList() selector.Voices {
	v := selector.Voices{}
	v.Append("сыендук", "syenduk")
	v.Append("кубик", "кубе", "kubik", "kube")
	v.Append("кураж", "бомбей", "kurazh", "bombej")
	v.Append("lostfilm", "lost")
	v.Append("newstudio")
	v.Append("амедиа", "amedia")
	return v
}

func main() {
	log.Infof("media-station v%s", Version)
	defer log.Infof("DONE")

	configPath := flag.String("config", "/etc/media-station/media-station.yml", "Path to YAML configuration file")
	verbose := flag.Bool("verbose", false, "Enable extra logs")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Load configuration failed: %s", err)
	}
	log.Debugf("Config: %+v", conf)

	torrentService, err := torrents.New(conf.Storage)
	if err != nil {
		log.Fatalf("Start torrent service failed: %s", err)
	}
	defer torrentService.Stop()

	server := web.Server{
		DiscoveryService: discovery.NewService(conf.Discovery),
		TorrentService:   torrentService,
		Selector: selector.MovieSelector{
			MinSeasonSizeMB:     1024,
			MaxSeasonSizeMB:     50 * 1024,
			MinSeedersThreshold: 50,
			QualityPrior:        []string{"1080p", "720p", "480p"},
			VoiceList:           getVoicePriorityList(),
		},
	}
	if err = server.Run(conf.Http.Host, conf.Http.Port); err != nil {
		log.Fatalf("Run web server failed: %s", err)
	}
}
