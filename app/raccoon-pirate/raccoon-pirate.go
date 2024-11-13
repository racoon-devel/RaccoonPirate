package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/apex/log"
	"github.com/blang/semver"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
	"github.com/racoon-devel/raccoon-pirate/internal/discovery"
	"github.com/racoon-devel/raccoon-pirate/internal/representation"
	"github.com/racoon-devel/raccoon-pirate/internal/selector"
	"github.com/racoon-devel/raccoon-pirate/internal/torrents"
	"github.com/racoon-devel/raccoon-pirate/internal/web"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var Version = "0.0.0"

func main() {
	log.Infof("raccoon-pirate %s [ %s ]", Version, runtime.GOARCH)
	defer log.Info("DONE")

	configPath := flag.String("config", "/etc/raccoon-pirate/raccoon-pirate.yml", "Path to YAML configuration file")
	verbose := flag.Bool("verbose", false, "Enable extra logs")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
		selfupdate.EnableLog()
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Load configuration failed: %s", err)
	}
	log.Infof("Config: %+v", conf)

	if conf.Application.AutoUpdate {
		if err = doSelfUpdate(); err != nil {
			log.Warnf("Auto update failed: %s", err)
		}
	}

	dbase, err := db.Open(conf.Storage)
	if err != nil {
		log.Fatalf("Open database failed: %s", err)
	}
	defer dbase.Close()

	printRegisteredTorrents(dbase)

	reprService := representation.New(conf.Representation)
	defer reprService.Clean()

	torrentService, err := torrents.New(conf.Storage, dbase, reprService)
	if err != nil {
		log.Fatalf("Start torrent service failed: %s", err)
	}
	defer torrentService.Stop()

	discoveryService := discovery.NewService(conf.Discovery)

	if conf.Frontend.Http.Enabled {
		server := web.Server{
			DiscoveryService: discoveryService,
			TorrentService:   torrentService,
			SelectCriterion:  conf.Selector.GetCriterion(),
			Selector: selector.New(selector.Settings{
				MinSeasonSizeMB:     int64(conf.Selector.MinSeasonSize),
				MaxSeasonSizeMB:     int64(conf.Selector.MaxSeasonSize),
				MinSeedersThreshold: int64(conf.Selector.MinSeedersThreshold),
				QualityPrior:        conf.Selector.Quality,
				VoiceList:           selector.Voices(conf.Selector.Voices),
			}),
		}

		if err = server.Run(conf.Frontend.Http.Host, conf.Frontend.Http.Port); err != nil {
			log.Errorf("Run web server failed: %s", err)
		} else {
			defer server.Shutdown()
		}
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh
	log.Info("Shutdowning")
}

func printRegisteredTorrents(dbase db.Database) {
	out := "Registered torrents:\n"
	list, err := dbase.LoadAllTorrents()
	if err != nil {
		log.Fatalf("Retrieve torrents list failed: %s", err)
	}
	for _, t := range list {
		out += fmt.Sprintf("ID: %s, Type: %d, Title: '%s', BelongsTo: '%s'\n", t.ID, t.Type, t.Title, t.BelongsTo)
	}
	log.Info(out)
}

func doSelfUpdate() error {
	if Version == "0.0.0" {
		return nil
	}
	v, err := semver.Parse(Version[1:])
	if err != nil {
		return err
	}
	latest, err := selfupdate.UpdateSelf(v, "racoon-devel/RaccoonPirate")
	if err != nil {
		return err
	}

	if latest.Version.Equals(v) {
		log.Info("Nothing to update")
		return nil
	}
	log.Infof("Successfully updated to %s", latest.Version)
	return nil
}
