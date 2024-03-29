package starq

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	yaml "gopkg.in/yaml.v3"
)

func Run(opts Opts) error {
	fmt.Printf("Parsed opts: %+v\n", opts)

	// count the total number of rules
	nAddRules := len(opts.PrependRules) + len(opts.AppendRules)

	// create the prepend rules starting from 1
	prependRules, err := ParseRules(1, opts.PrependRules)
	if err != nil {
		return fmt.Errorf("failed to parse rules: %w", err)
	}

	// parse a config file (if provided)
	var configRules []Rule = nil
	var rules []Rule
	if len(opts.ConfigFile) > 0 {
		bytes, err := os.ReadFile(opts.ConfigFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		fmt.Printf("Loaded config:\n%s\n", string(bytes))

		var config Config
		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			return fmt.Errorf("failed to unmarshall config object: %w", err)
		}
		fmt.Printf("Parsed config: %+v\n", config)

		// allocate space for the rules from the config
		nAddRules += len(config.Rules)
		configRules = config.Rules
		rules = make([]Rule, 0, nAddRules)
	} else {
		rules = make([]Rule, 0, nAddRules)
	}

	// create the append rules starting from on the number of rules so far + 1
	appendRules, err := ParseRules(nAddRules+1, opts.AppendRules)
	if err != nil {
		return fmt.Errorf("failed to parse rules: %w", err)
	}

	// assign all the rules
	rules = append(rules, prependRules...)
	rules = append(rules, configRules...)
	rules = append(rules, appendRules...)

	fmt.Printf("Rules: %+v\n", rules)
	return ApplyRules(rules, opts.Input, opts.Output, opts.Errors)
}

func ParseRules(countFrom int, raw []string) ([]Rule, error) {
	rules := make([]Rule, len(raw))
	for i, r := range raw {
		rule, err := ParseRule(fmt.Sprintf("rule-%d", countFrom+i), r)
		if err != nil {
			return nil, err
		}
		rules[i] = rule
	}
	return rules, nil
}

func ParseRule(name string, raw string) (Rule, error) {
	return Rule{
		Name: name,
		Jq:   raw,
	}, nil
}

func ApplyRules(rules []Rule, input io.Reader, output io.Writer, errors io.Writer) error {
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
