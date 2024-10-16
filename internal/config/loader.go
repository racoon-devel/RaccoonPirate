package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
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
	return result, err
}
