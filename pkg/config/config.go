package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config file path cannot be empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("invalid YAML format: %w", err)
	}

	if len(config.Hosts) == 0 {
		return nil, fmt.Errorf("no hosts found in config file")
	}

	for i, host := range config.Hosts {
		if err := validateHost(host); err != nil {
			return nil, fmt.Errorf("invalid host at index %d: %w", i, err)
		}
	}

	return &config, nil
}

// ParseDomainsFromString parses a comma-separated string of domains
func ParseDomainsFromString(domains string) ([]string, error) {
	if domains == "" {
		return nil, fmt.Errorf("domains string cannot be empty")
	}

	parts := strings.Split(domains, ",")
	var hosts []string

	for i, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		if err := validateHost(trimmed); err != nil {
			return nil, fmt.Errorf("invalid domain at position %d (%s): %w", i+1, trimmed, err)
		}

		hosts = append(hosts, trimmed)
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("no valid domains found in the provided string")
	}

	return hosts, nil
}

// validateHost validates a host string format
func validateHost(host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Handle IPv6 addresses with brackets [::1]:8080
	if strings.HasPrefix(host, "[") {
		closeBracket := strings.Index(host, "]")
		if closeBracket == -1 {
			return fmt.Errorf("invalid IPv6 address format (missing closing bracket): %s", host)
		}

		ipv6Addr := host[1:closeBracket]
		if ipv6Addr == "" {
			return fmt.Errorf("IPv6 address cannot be empty")
		}

		remainder := host[closeBracket+1:]
		if remainder != "" && !strings.HasPrefix(remainder, ":") {
			return fmt.Errorf("invalid format after IPv6 address: %s", host)
		}

		return nil
	}

	parts := strings.Split(host, ":")

	// If more than 2 parts, could be IPv6 without brackets or invalid format
	if len(parts) > 2 {
		// Check if it looks like an IPv6 address (contains multiple colons)
		if strings.Count(host, ":") > 1 {
			// Assume it's IPv6 without brackets - this is valid
			return nil
		}
		return fmt.Errorf("invalid host format (too many colons): %s", host)
	}

	hostname := strings.TrimSpace(parts[0])
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Basic hostname validation - could be enhanced with regex
	if strings.Contains(hostname, " ") {
		return fmt.Errorf("hostname cannot contain spaces: %s", hostname)
	}

	// Validate port if present
	if len(parts) == 2 {
		portStr := strings.TrimSpace(parts[1])
		if portStr != "" {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return fmt.Errorf("invalid port number: %s", portStr)
			}
			if port <= 0 || port > 65535 {
				return fmt.Errorf("port number out of range (1-65535): %d", port)
			}
		}
	}

	return nil
}

// Validate validates the application configuration
func (c *AppConfig) Validate() error {
	if c.ConfigFile == "" && c.Domains == "" {
		return fmt.Errorf("either --config or --domains must be specified")
	}

	if c.ConfigFile != "" && c.Domains != "" {
		return fmt.Errorf("--config and --domains cannot be used together")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.OutputFormat != "" && c.OutputFormat != "table" && c.OutputFormat != "json" && c.OutputFormat != "yaml" {
		return fmt.Errorf("invalid output format: %s (supported: table, json, yaml)", c.OutputFormat)
	}

	return nil
}

// GetHosts returns the list of hosts based on the configuration
func (c *AppConfig) GetHosts() ([]string, error) {
	if c.ConfigFile != "" {
		config, err := LoadConfig(c.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
		return config.Hosts, nil
	}

	if c.Domains != "" {
		hosts, err := ParseDomainsFromString(c.Domains)
		if err != nil {
			return nil, fmt.Errorf("failed to parse domains: %w", err)
		}
		return hosts, nil
	}

	return nil, fmt.Errorf("no hosts configuration provided")
}
