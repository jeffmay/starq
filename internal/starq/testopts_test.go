package starq_test

import (
	"starq/internal/starq"
)

// TestOpts replaces the [starq.Opts] struct and implements the [starq.TransformerLoader] interface.
//
// The big difference is that this stores [TestTransformer]s instead of [starq.Transformer]s and
// it loads these transformers from file upon construction, rather than lazily. You should not use
// this if you intend to test the [starq.Opts].
//
// This is mainly useful when you want to stub the input and output of the transformers, rather
// than needing to read the output from the filesystem or stdout and for stubbing the input with
// strings instead of needing to use files. This is primarily used for testing the [starq.Runner].
type TestOpts struct {
	starq.GlobalConfig
	Transformers []TestTransformer
}

var _ starq.TransformerLoader = new(TestOpts)

func (o TestOpts) LoadTransformers() ([]starq.Transformer, error) {
	transformers := make([]starq.Transformer, len(o.Transformers))
	for i, t := range o.Transformers {
		transformers[i] = t
	}
	return transformers, nil
}

func (o TestOpts) GetGlobalConfig() starq.GlobalConfig {
	return o.GlobalConfig
}

// MakeTestOpts uses the builder pattern to construct a new [TestOpts] with stubbed [TestTransformer]s.
//
// By default, the [TestOpts] will use the [DefaultTestTransformer] as its only transformer.
func MakeTestOpts() *TestOpts {
	return &TestOpts{
		GlobalConfig: starq.GlobalConfig{
			PrependRules: nil,
			AppendRules:  nil,
		},
		Transformers: []TestTransformer{*DefaultTestTransformer()},
	}
}

func (o *TestOpts) SetAppendRules(rules ...string) *TestOpts {
	o.GlobalConfig.AppendRules = rules
	return o
}

func (o *TestOpts) WithRulesAppended(rules ...string) *TestOpts {
	o.GlobalConfig.AppendRules = append(o.GlobalConfig.AppendRules, rules...)
	return o
}

func (o *TestOpts) SetPrependRules(rules ...string) *TestOpts {
	o.GlobalConfig.PrependRules = rules
	return o
}

func (o *TestOpts) WithRulesPrepended(rules ...string) *TestOpts {
	o.GlobalConfig.PrependRules = append(o.GlobalConfig.PrependRules, rules...)
	return o
}

func (o *TestOpts) WithTransformers(transformers ...*TestTransformer) *TestOpts {
	unwrapped := make([]TestTransformer, 0, len(transformers))
	for _, c := range transformers {
		if c != nil {
			unwrapped = append(unwrapped, *c)
		}
	}
	o.Transformers = unwrapped
	return o
}
