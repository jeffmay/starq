package starq

import (
	"fmt"
	"io"
	"os/exec"
	"starq/internal/iox"
)

// runner is a dependency-injected runner for the starq CLI.
//
// The runner can run every transformer provided via a [TransformerLoader], or just a single [Transformer].
//
// All of the methods eventually invoke the [ExecuteRules] function with the appropriate arguments.
type runner struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewRunner constructs a new runner with the given input/output/error streams.
//
// This is the public constructor for the runner, since you should never provide nil
// arguments for any of the streams.
func NewRunner(stdin io.Reader, stdout io.Writer, stderr io.Writer) *runner {
	if stdin == nil || stdout == nil || stderr == nil {
		panic("nil streams are not allowed")
	}
	return &runner{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// TransformerLoader loads all transformers from configs and provides the [GlobalConfig].
type TransformerLoader interface {
	GetGlobalConfig() GlobalConfig
	LoadTransformers() ([]Transformer, error)
}

// GlobalConfig is a struct that holds configuration that is used across all [Transformer]s.
type GlobalConfig struct {
	PrependRules []string
	AppendRules  []string
}

// ExecuteRules runs `jq` command(s) with the given rules and input/output/error streams.
func ExecuteRules(rules []Rule, input io.Reader, output io.Writer, errors io.Writer) error {
	// apply all rules in sequence in a single command
	ruleStrings := make([]string, len(rules))
	for i, r := range rules {
		ruleStrings[i] = r.Jq
	}
	cmd := exec.Command("jq", ruleStrings...)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = errors

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run jq: %w", err)
	}
	return nil
}

// CombineRules appends & prepents the rules from the [GlobalConfig] appropriately.
func CombineRules(rules []Rule, global GlobalConfig) ([]Rule, error) {
	// count the total number of rules
	nAddRules := len(global.PrependRules) + len(global.AppendRules)

	// create the prepend rules starting from 1
	prependRules := NameJqRules(1, global.PrependRules)

	// allocate space for the rules from the config
	nConfigRules := len(rules) + nAddRules

	// create the append rules starting from on the number of rules so far + 1
	appendRules := NameJqRules(nConfigRules+1, global.AppendRules)

	// fill out the rules
	allRules := make([]Rule, 0, nConfigRules)
	allRules = append(allRules, prependRules...)
	allRules = append(allRules, rules...)
	allRules = append(allRules, appendRules...)

	return allRules, nil
}

// RunTransformer runs a single [Transformer] with the given [GlobalConfig].
func (r *runner) RunTransformer(transformer Transformer, global GlobalConfig) error {
	input, err := transformer.InputReader(r.stdin)
	if err != nil {
		return err
	}
	if input == nil {
		input = io.NopCloser(r.stdin)
	} else {
		defer input.Close()
	}
	defn := transformer.Config()
	rules, err := CombineRules(defn.Rules, global)
	if err != nil {
		return err
	}
	output, err := transformer.OutputWriter(r.stdout)
	if err != nil {
		return err
	}
	if output == nil {
		output = iox.NopWriterCloser(r.stdout)
	} else {
		defer output.Close()
	}
	errors, err := transformer.ErrorsWriter(r.stderr)
	if err != nil {
		return err
	}
	if errors == nil {
		errors = iox.NopWriterCloser(r.stderr)
	} else {
		defer errors.Close()
	}
	return ExecuteRules(rules, input, output, errors)
}

// RunAllTransformers runs all [Transformer]s loaded from the given [TransformerLoader] along with the [GlobalConfig].
func (r *runner) RunAllTransformers(loader TransformerLoader) error {
	transformers, err := loader.LoadTransformers()
	if err != nil {
		return fmt.Errorf("failed to load transformers: %w", err)
	}
	for _, transformer := range transformers {
		err := r.RunTransformer(transformer, loader.GetGlobalConfig())
		if err != nil {
			return err
		}
	}
	return nil
}
