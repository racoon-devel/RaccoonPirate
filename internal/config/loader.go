package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type schemaContext struct {
	source  *string
	schema  *jsonschema.Schema
	migrate migrateFn
}

const currentConfigVersion = 2

//go:embed config-schema.v1.json
var schemaSourceV1 string

//go:embed config-schema.v2.json
var schemaSourceV2 string

var schemas = map[int]*schemaContext{
	1: {source: &schemaSourceV1, migrate: migrateV1ToV2},
	2: {source: &schemaSourceV2},
}

type versionOnlyConfig struct {
	Application struct {
		ConfigVersion uint `json:"config-version"`
	}
}

func init() {
	for _, schemaCtx := range schemas {
		schemaJs, err := jsonschema.UnmarshalJSON(strings.NewReader(*schemaCtx.source))
		if err != nil {
			panic(fmt.Sprintf("load embedded schema failed: %s", err))
		}
		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource("config-schema.json", schemaJs); err != nil {
			panic(fmt.Sprintf("load embedded schema failed: %s", err))
		}
		schemaCtx.schema = compiler.MustCompile("config-schema.json")
	}
}

func Load(destination string) (Config, error) {
	schema := schemas[currentConfigVersion]

	content, err := os.ReadFile(destination)
	if err != nil {
		return Config{}, err
	}

	jsonRaw, err := yaml.YAMLToJSON(content)
	if err != nil {
		return Config{}, err
	}

	var versionProber versionOnlyConfig
	_ = json.Unmarshal(jsonRaw, &versionProber)

	version := &versionProber.Application.ConfigVersion
	if *version != currentConfigVersion {
		if *version == 0 {
			*version = 1
		}

		log.Warnf("Deprecated version of config detected: %d, trying to migrate...", *version)
		jsonRaw, err = migrateConfig(*version, jsonRaw)
		if err != nil {
			return Config{}, fmt.Errorf("migrate config failed: %w", err)
		}
	}

	j, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonRaw))
	if err != nil {
		return Config{}, err
	}

	err = schema.schema.Validate(j)
	if err != nil {
		return Config{}, err
	}

	var result Config
	err = json.Unmarshal(jsonRaw, &result)
	if err != nil {
		return Config{}, err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("get user home directory failed: %w", err)
	}

	if !filepath.IsAbs(result.Storage.Directory) {
		result.Storage.Directory = filepath.Join(userHomeDir, result.Storage.Directory)
	}

	if !filepath.IsAbs(result.Representation.Directory) {
		result.Representation.Directory = filepath.Join(userHomeDir, result.Representation.Directory)
	}

	return result, err
}
