package starq

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type TransformerConfig struct {
	Source ConfigSource `yaml:"-"`
	Input  *Input       `yaml:"input"`
	Rules  []Rule       `yaml:"rules"`
}

type Input struct {
	File *string `yaml:"file"`
}

type Rule struct {
	Name string `yaml:"name"`
	Jq   string `yaml:"jq"`
}

// DefaultTransformerConfig returns a new [TransformerConfig] with the [DefaultConfigSource] and no rules or input file.
//
// This means that all input comes from stdin and the only rules applied come from the [starq.GlobalConfig].
func DefaultTransformerConfig() *TransformerConfig {
	return &TransformerConfig{
		Source: DefaultConfigSource,
		Input:  nil,
		Rules:  make([]Rule, 0),
	}
}

type ConfigSource string

// DefaultConfigSource is the source for configuration when the config is not loaded from a file.
const DefaultConfigSource = ConfigSource("default")

// FileConfigSource returns a [ConfigSource] for a file.
func FileConfigSource(filename string) (ConfigSource, error) {
	var empty ConfigSource
	abspath, err := filepath.Abs(strings.TrimPrefix(filename, "file://"))
	if err != nil {
		return empty, fmt.Errorf("failed to resolve absolute path to config file: %w", err)
	}
	return ConfigSource(fmt.Sprintf("file://%s", abspath)), nil
}

func (src ConfigSource) String() string {
	return string(src)
}

// AsFileName returns the source name or file name and a boolean indicating whether it is a file name.
func (src ConfigSource) AsFileName() (string, bool) {
	s := string(src)
	path := strings.TrimPrefix(s, "file://")
	// if the path was trimmed, then it was a file source
	return path, path != s
}

// ReadTransformerConfigFile reads a [TransformerConfig] from a file.
func ReadTransformerConfigFile(filename string) (TransformerConfig, error) {
	var config TransformerConfig
	source, err := FileConfigSource(filename)
	if err != nil {
		return config, err
	}
	absPath, _ := source.AsFileName()
	bytes, err := os.ReadFile(absPath)
	if err != nil {
		return config, fmt.Errorf("failed to read file '%s': %w", absPath, err)
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse transformer config file: %w", err)
	}
	config.Source = source
	return config, err
}
