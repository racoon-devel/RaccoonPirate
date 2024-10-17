package main

import (
	"flag"

	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
	"github.com/racoon-devel/raccoon-pirate/internal/discovery"
	"github.com/racoon-devel/raccoon-pirate/internal/selector"
	"github.com/racoon-devel/raccoon-pirate/internal/torrents"
	"github.com/racoon-devel/raccoon-pirate/internal/web"
)

var Version = "0.0.0"

func main() {
	log.Infof("raccoon-pirate %s", Version)
	defer log.Infof("DONE")

	configPath := flag.String("config", "/etc/raccoon-pirate/raccoon-pirate.yml", "Path to YAML configuration file")
	verbose := flag.Bool("verbose", false, "Enable extra logs")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Load configuration failed: %s", err)
	}
	log.Infof("Config: %+v", conf)

	dbase, err := db.Open(conf.Storage)
	if err != nil {
		log.Fatalf("Open database failed: %s", err)
	}
	defer dbase.Close()

	torrentService, err := torrents.New(conf.Storage)
	if err != nil {
		log.Fatalf("Start torrent service failed: %s", err)
	}
	defer torrentService.Stop()

	if conf.Frontend.Http.Enabled {
		server := web.Server{
			DiscoveryService: discovery.NewService(conf.Discovery),
			TorrentService:   torrentService,
			Selector: selector.MovieSelector{
				MinSeasonSizeMB:     int64(conf.Selector.MinSeasonSize),
				MaxSeasonSizeMB:     int64(conf.Selector.MaxSeasonSize),
				MinSeedersThreshold: int64(conf.Selector.MinSeedersThreshold),
				QualityPrior:        conf.Selector.Quality,
				VoiceList:           selector.Voices(conf.Selector.Voices),
			},
		}

		if err = server.Run(conf.Frontend.Http.Host, conf.Frontend.Http.Port); err != nil {
			log.Fatalf("Run web server failed: %s", err)
		}
	}
}
