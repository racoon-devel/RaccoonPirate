package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:embed config-schema.json
var schemaSource string

var schemaCompiler *jsonschema.Compiler
var schema *jsonschema.Schema

func init() {
	schemaJs, err := jsonschema.UnmarshalJSON(strings.NewReader(schemaSource))
	if err != nil {
		panic(fmt.Sprintf("load embedded schema failed: %s", err))
	}

	schemaCompiler = jsonschema.NewCompiler()
	if err := schemaCompiler.AddResource("config-schema.json", schemaJs); err != nil {
		panic(fmt.Sprintf("load embedded schema failed: %s", err))
	}
	schema = schemaCompiler.MustCompile("config-schema.json")
}

func Load(destination string) (Config, error) {

	content, err := os.ReadFile(destination)
	if err != nil {
		return Config{}, err
	}

	jsonRaw, err := yaml.YAMLToJSON(content)
	if err != nil {
		return Config{}, err
	}

	j, err := jsonschema.UnmarshalJSON(bytes.NewReader(jsonRaw))
	if err != nil {
		return Config{}, err
	}

	err = schema.Validate(j)
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
