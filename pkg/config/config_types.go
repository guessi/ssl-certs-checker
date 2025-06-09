package config

type Config struct {
	Hosts []string `yaml:"hosts"`
}

type AppConfig struct {
	ConfigFile   string
	Domains      string
	Timeout      int
	Insecure     bool
	OutputFormat string
}
