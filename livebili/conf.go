package livebili

type Config struct {
	Uids          []int64 `yaml:"uids" mapstructure:"uids"`
	CheckDuration int     `yaml:"check_duration" mapstructure:"check_duration"`
	Groups        []int64 `yaml:"groups" mapstructure:"groups"`
}
