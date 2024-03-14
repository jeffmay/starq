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
	PETSTORE_TO_READONLY_STDOUT_PATH RelPath = "file://./config/petstore-to-readonly-stdout.yaml"
	PETSTORE_TO_YAML_STDOUT_PATH     RelPath = "file://./config/petstore-to-yaml-stdout.yaml"
	SIMPLE_PATH                      RelPath = "file://./config/simple.yaml"

	// Input files:
	PETSTORE_OPENAPI_JSON_PATH RelPath = "file://./input/petstore-openapi.json"
)

//go:embed input/petstore-openapi.json
var PETSTORE_OPENAPI_JSON string
