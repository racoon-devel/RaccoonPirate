package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type migrateFn func([]byte) ([]byte, error)

func migrateConfig(version uint, jsonRaw []byte) ([]byte, error) {
	if version < 1 || version >= currentConfigVersion {
		return nil, fmt.Errorf("unknown config version: %d", version)
	}

	var err error
	for nextVersion := version; nextVersion < currentConfigVersion; nextVersion++ {
		migrator := schemas[int(nextVersion)].migrate
		jsonRaw, err = migrator(jsonRaw)
		if err != nil {
			return nil, fmt.Errorf("migration from %d to %d failed: %w", nextVersion, nextVersion+1, err)
		}
	}

	return jsonRaw, nil
}

type storageV1 struct {
	Directory   string
	Driver      string
	Limit       uint
	AddTimeout  uint `json:"add-timeout"`
	ReadTimeout uint `json:"read-timeout"`
	TTL         uint
}

type configV1 struct {
	Frontend       Frontend
	Application    Application
	Api            Api
	Discovery      Discovery
	Storage        storageV1
	Representation Representation
	Selector       Selector
}

func migrateV1ToV2(jsonRaw []byte) ([]byte, error) {
	var old configV1
	if err := json.Unmarshal(jsonRaw, &old); err != nil {
		return nil, err
	}

	new := Config{
		Frontend:       old.Frontend,
		Application:    old.Application,
		Api:            old.Api,
		Discovery:      old.Discovery,
		Representation: old.Representation,
		Selector:       old.Selector,
		Database: Database{
			Path:   filepath.Join(old.Storage.Directory, "database.db"),
			Driver: old.Storage.Driver,
		},
		Torrent: Torrent{
			Driver: "builtin",
			Builtin: Builtin{
				Directory:   old.Storage.Directory,
				Limit:       old.Storage.Limit,
				AddTimeout:  old.Storage.AddTimeout,
				ReadTimeout: old.Storage.ReadTimeout,
				TTL:         old.Storage.TTL,
			},
			TorrServer: TorrServer{
				URL: "http://127.0.0.1:8090/",
			},
		},
	}

	new.Application.ConfigVersion = 2

	return json.Marshal(new)
}
