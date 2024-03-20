package sample

import (
	_ "embed"
	"path/filepath"
	"strings"
)

func MustNormalize(pathToProjectRoot string, pathInThisFile RelPath) string {
	absPath, err := Normalize(pathToProjectRoot, pathInThisFile)
	if err != nil {
		panic(err)
	}
	return absPath
}

func Normalize(pathToProjectRoot string, pathInThisFile RelPath) (string, error) {
	return filepath.Abs(filepath.Join(pathToProjectRoot, "sample", strings.TrimPrefix(string(pathInThisFile), "file://./")))
}

// RelPath is a path to a file in the project, relative to the sample package directory.
type RelPath string

// Paths:
// These all start with "file://" so that the IDE will render them as links, but this is stripped by the Normalize function.
const (
	// Config files:
	CONFIG_PETSTORE_JSON_TO_READONLY_PATH RelPath = "file://./config/petstore-json-to-readonly.yaml"
	CONFIG_PETSTORE_JSON_TO_YAML_PATH     RelPath = "file://./config/petstore-json-to-yaml.yaml"
	CONFIG_SIMPLE_PATH                    RelPath = "file://./config/simple.yaml"

	// Input files:
	INPUT_PETSTORE_OPENAPI_JSON_PATH RelPath = "file://./input/petstore-openapi.json"

	// Output files:
	OUTPUT_PETSTORE_OPENAPI_YAML_PATH          RelPath = "file://./output/petstore-openapi-readonly.json"
	OUTPUT_PETSTORE_OPENAPI_READONLY_JSON_PATH RelPath = "file://./output/petstore-openapi-readonly.json"
)

//go:embed input/petstore-openapi.json
var PETSTORE_OPENAPI_JSON string
