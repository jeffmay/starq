package starq

import "io"

type Opts struct {
	PrependRules []string
	AppendRules  []string
	ConfigFile   string
	Input        io.Reader
	Output       io.Writer
}
