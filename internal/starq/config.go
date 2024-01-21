package starq

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Name string `yaml:"name"`
	Jq   string `yaml:"jq"`
}
