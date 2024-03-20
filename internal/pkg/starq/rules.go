package starq

import "fmt"

// NameJqRules adds names like "rule-1", "rule-2", ... to a list of given `jq` expressions starting from some initial count.
func NameJqRules(nameCountFrom int, jqRules []string) []Rule {
	rules := make([]Rule, len(jqRules))
	for i, jq := range jqRules {
		rules[i] = Rule{Name: fmt.Sprintf("rule-%d", nameCountFrom+i), Jq: jq}
	}
	return rules
}
