package main

import (
	"flag"
	"github.com/apex/log"
	"github.com/racoon-devel/media-station/internal/config"
	"github.com/racoon-devel/media-station/internal/web"
)

var Version = "0.0.0"

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

	server := web.Server{}
	if err = server.Run(conf.Http.Host, conf.Http.Port); err != nil {
		log.Fatalf("Run web server failed: %s", err)
	}
}
