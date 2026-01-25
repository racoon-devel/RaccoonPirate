package main

import (
	"os"

	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: dbconv <input_path> <input_driver> <output_path> <output_driver>")
	}

	input := config.Database{
		Path:   os.Args[1],
		Driver: os.Args[2],
	}
	output := config.Database{
		Path:   os.Args[3],
		Driver: os.Args[4],
	}

	inputDb, err := db.Open(input)
	if err != nil {
		log.Fatalf("Open input database failed: %s", err)
	}
	defer inputDb.Close()

	outputDb, err := db.Open(output)
	if err != nil {
		log.Fatalf("Open output database failed: %s", err)
	}
	defer outputDb.Close()

	version, _ := inputDb.GetVersion()
	_ = outputDb.SetVersion(version)

	torrents, err := inputDb.LoadTorrents(true)
	if err != nil {
		log.Fatalf("Load torrents failed: %s", err)
	}

	for _, t := range torrents {
		if err := outputDb.PutTorrent(t); err != nil {
			log.Warnf("Add torrent %s failed: %s", t.ID, err)
		}
	}

	log.Info("Database converted")
}
