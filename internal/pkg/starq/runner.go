package starq

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/jeffmay/starq/internal/pkg/iox"

	"github.com/goccy/go-yaml"
)

// Streams abstracts the input/output/error streams, so that they can be passed around in a single object
type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func NewStreams(in io.Reader, out, err io.Writer) Streams {
	return Streams{
		Stdin:  in,
		Stdout: out,
		Stderr: err,
	}
}

// Attach will provide the input/output/error streams (if defined) to the given command subprocess streams.
func (s Streams) Attach(cmd *exec.Cmd) *exec.Cmd {
	if s.Stdin != nil {
		cmd.Stdin = s.Stdin
	}
	if s.Stdout != nil {
		cmd.Stdout = s.Stdout
	}
	if s.Stderr != nil {
		cmd.Stderr = s.Stderr
	}
	return cmd
}

// runner is a dependency-injected runner for the starq CLI.
//
// The runner can run every transformer provided via a [TransformerLoader], or just a single [Transformer].
//
// All of the methods eventually invoke the [ExecuteJq] function with the appropriate arguments.
type runner Streams

// NewRunner constructs a new runner with the given input/output/error streams.
//
// This is the public constructor for the runner, since you should never provide nil
// arguments for any of the streams.
func NewRunner(stdin io.Reader, stdout io.Writer, stderr io.Writer) *runner {
	if stdin == nil || stdout == nil || stderr == nil {
		panic("nil streams are not allowed")
	}
	return &runner{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
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

// ExecuteJq runs `jq` command(s) with the given rules and input/output/error streams.
//
// NOTE: This requires that `jq` is installed and available in the PATH running this program.
func ExecuteJq(rules []Rule, streams Streams) error {
	ruleStrings := make([]string, len(rules))
	for i, r := range rules {
		ruleStrings[i] = r.Jq
	}
	// apply all rules in sequence in a single command
	cmd := streams.Attach(exec.Command("jq", ruleStrings...))
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

// GuessFormatFromFileExt attempts to guess the [DocumentFormat] from the given file extension.
func GuessFormatFromFileExt(filename string) DocumentFormat {
	if strings.HasSuffix(filename, ".json") {
		return JSONFormat
	}
	if strings.HasSuffix(filename, ".yaml") {
		return YAMLFormat
	}
	return UnknownFormat
}

// RunTransformer runs a single [Transformer] with the given [GlobalConfig].
//
// TODO: Break this into smaller steps
func (r *runner) RunTransformer(transformer Transformer, global GlobalConfig) error {
	defn := transformer.Config()
	input, err := transformer.InputReader(r.Stdin)
	if err != nil {
		return err
	}
	if input == nil {
		input = io.NopCloser(r.Stdin)
	} else {
		defer input.Close()
	}
	rules, err := CombineRules(defn.Rules, global)
	if err != nil {
		return err
	}
	output, err := transformer.OutputWriter(r.Stdout)
	if err != nil {
		return err
	}
	if output == nil {
		output = iox.NopWriterCloser(r.Stdout)
	} else {
		defer output.Close()
	}
	errors, err := transformer.ErrorsWriter(r.Stderr)
	if err != nil {
		return err
	}
	if errors == nil {
		errors = iox.NopWriterCloser(r.Stderr)
	} else {
		defer errors.Close()
	}
	// convert YAML input to JSON input and JSON output to YAML output, as needed
	// set the input format
	inFormat := UnknownFormat
	inFile := ""
	if defn.Input != nil {
		inFormat = defn.Input.Format
		inFile = defn.Input.File
	}
	if inFormat == UnknownFormat {
		// attempt to guess the format from the input file extension
		inFormat = GuessFormatFromFileExt(inFile)
	}
	if inFormat == UnknownFormat {
		// if we still don't know, warn about it and attempt YAML
		// TODO: Warn about using YAML as a default
		inFormat = YAMLFormat
	}
	// set the output format
	outFormat := UnknownFormat
	if defn.Output != nil {
		outFormat = defn.Output.Format
	}
	if outFormat == UnknownFormat {
		// default to the input format when the output format is not specified
		outFormat = inFormat
	}
	if inFormat == YAMLFormat {
		// attempt to parse as YAML
		inBytes, err := io.ReadAll(input)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		jsonBytes, err := yaml.YAMLToJSON(inBytes)
		if err != nil {
			return fmt.Errorf("failed to convert input to JSON: %w", err)
		}
		// replace the input stream with the JSON input
		input = io.NopCloser(strings.NewReader(string(jsonBytes)))
	}
	// set the output format
	var jsonOut io.Writer = output
	var capturedJsonOut *strings.Builder = nil
	if outFormat == YAMLFormat {
		// replace the output stream with a string builder that captures the JSON output
		capturedJsonOut = new(strings.Builder)
		jsonOut = capturedJsonOut
	}
	// run the jq command
	err = ExecuteJq(rules, NewStreams(input, jsonOut, errors))
	if err != nil {
		return err
	}
	// write the final YAML output
	if outFormat == YAMLFormat {
		// convert the JSON output to YAML
		jsonBytes, err := io.ReadAll(strings.NewReader(capturedJsonOut.String()))
		if err != nil {
			return fmt.Errorf("failed to read JSON output: %w", err)
		}
		yamlBytes, err := yaml.JSONToYAML(jsonBytes)
		if err != nil {
			return fmt.Errorf("failed to convert JSON to YAML: %w", err)
		}
		// write the final output as YAML
		_, err = output.Write(yamlBytes)
		if err != nil {
			return fmt.Errorf("failed to write YAML output: %w", err)
		}
	}
	return nil
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
