package starq_test

import (
	"encoding/json"
	"fmt"
	"io"
	"starq/internal/starq"
	"strings"

	"github.com/valyala/fastjson"
)

type TestOpts struct {
	PrependRules []string
	AppendRules  []string
	ConfigFile   string
	Input        io.Reader
	output       *strings.Builder
	errors       *strings.Builder
}

func (o TestOpts) Output() string {
	return o.output.String()
}

func (o TestOpts) MustOutputValue() *fastjson.Value {
	return fastjson.MustParse(o.Output())
}

func (o TestOpts) Pretty() string {
	var obj interface{}
	err := json.Unmarshal([]byte(o.Output()), &obj)
	if err != nil {
		panic(fmt.Errorf("could not indent JSON: %w", err))
	}
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not indent JSON: %w", err))
	}
	return string(out)
}

func (o TestOpts) Errors() string {
	return o.errors.String()
}

func (o TestOpts) Opts() starq.Opts {
	return starq.Opts{
		PrependRules: o.PrependRules,
		AppendRules:  o.AppendRules,
		ConfigFile:   o.ConfigFile,
		Input:        o.Input,
		Output:       o.output,
		Errors:       o.errors,
	}
}

type TestOpt func(*TestOpts)

func WithAppendRules(rules ...string) TestOpt {
	return func(o *TestOpts) {
		o.AppendRules = rules
	}
}

func AppendRules(rules ...string) TestOpt {
	return func(o *TestOpts) {
		o.AppendRules = append(o.AppendRules, rules...)
	}
}

func WithPrependRules(rules ...string) TestOpt {
	return func(o *TestOpts) {
		o.PrependRules = rules
	}
}

func PrependRules(rules ...string) TestOpt {
	return func(o *TestOpts) {
		o.PrependRules = append(o.PrependRules, rules...)
	}
}

func WithConfigFile(configFile string) TestOpt {
	return func(o *TestOpts) {
		o.ConfigFile = configFile
	}
}

func WithInputString(input string) TestOpt {
	return func(o *TestOpts) {
		reader := strings.NewReader(input)
		o.Input = reader
	}
}

func WithInputReader(input io.Reader) TestOpt {
	return func(o *TestOpts) {
		o.Input = input
	}
}

func MakeTestOpts(opts ...TestOpt) TestOpts {
	var o TestOpts
	for _, opt := range opts {
		opt(&o)
	}
	if o.Input == nil {
		o.Input = strings.NewReader("")
	}
	o.output = new(strings.Builder)
	o.errors = new(strings.Builder)
	return o
}
