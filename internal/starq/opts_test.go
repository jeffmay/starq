package starq_test

import (
	"starq/internal/starq"
	"strings"
)

type Opt func(*starq.Opts)

func WithInputString(input string) Opt {
	return func(o *starq.Opts) {
		reader := strings.NewReader(input)
		o.Input = reader
	}
}

func WithOutputWriter(output *strings.Builder) Opt {
	return func(o *starq.Opts) {
		o.Output = output
	}
}

func WithConfigFile(configFile string) Opt {
	return func(o *starq.Opts) {
		o.ConfigFile = configFile
	}
}

func MakeOpts(opts ...Opt) starq.Opts {
	var o starq.Opts
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
