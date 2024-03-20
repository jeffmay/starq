package starq

// Opts is a [TransformerLoader] that is based on config files and global config --
// typically provided via CLI arguments.
type Opts struct {
	GlobalConfig
	ConfigFiles []string
}

var _ TransformerLoader = new(Opts)

func (opts Opts) GetGlobalConfig() GlobalConfig {
	return opts.GlobalConfig
}

// IsEmpty returns true if no config files or rules are set.
func (opts Opts) IsEmpty() bool {
	return len(opts.PrependRules) == 0 && len(opts.AppendRules) == 0 && len(opts.ConfigFiles) == 0
}

// LoadTransformers loads all the transformer config files into memory.
func (opts Opts) LoadTransformers() ([]Transformer, error) {
	configs := make([]Transformer, len(opts.ConfigFiles))
	for i, filename := range opts.ConfigFiles {
		config, err := ReadTransformerConfigFile(filename)
		if err != nil {
			return nil, err
		}
		configs[i] = &config
	}
	return configs, nil
}
