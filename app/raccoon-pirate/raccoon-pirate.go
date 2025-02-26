package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/RacoonMediaServer/rms-library/pkg/selector"
	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
	"github.com/racoon-devel/raccoon-pirate/internal/discovery"
	"github.com/racoon-devel/raccoon-pirate/internal/frontend"
	"github.com/racoon-devel/raccoon-pirate/internal/remote"
	"github.com/racoon-devel/raccoon-pirate/internal/representation"
	"github.com/racoon-devel/raccoon-pirate/internal/smartsearch"
	"github.com/racoon-devel/raccoon-pirate/internal/telegram"
	"github.com/racoon-devel/raccoon-pirate/internal/torrents"
	"github.com/racoon-devel/raccoon-pirate/internal/updater"
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

	dbase, err := db.Open(conf.Storage)
	if err != nil {
		log.Fatalf("Open database failed: %s", err)
	}
	defer dbase.Close()

	if conf.Application.AutoUpdate {
		u := updater.Updater{
			CurrentVersion: Version,
			Storage:        dbase,
		}

		updated, err := u.TryUpdate()
		if err != nil {
			log.Warnf("Auto update failed: %s", err)
		}
		if updated {
			dbase.Close()
			if err = u.Restart(); err != nil {
				log.Fatalf("Restart application after update failed: %s", err)
			}
		}
	}

	printRegisteredTorrents(dbase)

	apiConn := &remote.Connector{Config: conf.Api}
	if err = apiConn.ObtainToken(); err != nil {
		log.Errorf("!!! Discovery and Telegram services will not work, because obtain API token failed: %s !!!", err)
	}

	reprService := representation.New(conf.Representation)
	defer reprService.Clean()

	torrentService, err := torrents.New(conf.Storage, dbase, reprService)
	if err != nil {
		log.Fatalf("Start torrent service failed: %s", err)
	}
	defer torrentService.Stop()

	discoveryService := discovery.NewService(apiConn, conf.Discovery)
	smartSearchService := smartsearch.NewService(apiConn, conf.Discovery)

	frontendSetup := frontend.Setup{
		DiscoveryService:   discoveryService,
		TorrentService:     torrentService,
		SmartSearchService: smartSearchService,
		SelectCriterion:    conf.Selector.GetCriterion(),
		Selector: selector.New(selector.Settings{
			MinSeasonSizeMB:     int64(conf.Selector.MinSeasonSize),
			MaxSeasonSizeMB:     int64(conf.Selector.MaxSeasonSize),
			MinSeedersThreshold: int64(conf.Selector.MinSeedersThreshold),
			QualityPrior:        conf.Selector.Quality,
			VoiceList:           selector.Voices(conf.Selector.Voices),
		}),
	}

	if conf.Frontend.Http.Enabled {
		webServer := web.Server{Setup: frontendSetup}
		if err = webServer.Run(conf.Frontend.Http.Host, conf.Frontend.Http.Port); err != nil {
			log.Errorf("Run web server failed: %s", err)
		} else {
			defer webServer.Shutdown()
		}
	}

	if conf.Frontend.Telegram.Enabled {
		bot := telegram.Bot{
			Setup:     frontendSetup,
			Transport: apiConn.NewBotSession(conf.Frontend.Telegram.ApiPath),
		}
		bot.Run()
		defer bot.Shutdown()
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
