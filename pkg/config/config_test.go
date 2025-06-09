package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDomainsFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "single domain",
			input: "example.com",
			want:  []string{"example.com"},
		},
		{
			name:  "multiple domains",
			input: "example.com,google.com,github.com",
			want:  []string{"example.com", "google.com", "github.com"},
		},
		{
			name:  "domains with ports",
			input: "example.com:443,google.com:8080",
			want:  []string{"example.com:443", "google.com:8080"},
		},
		{
			name:  "domains with spaces",
			input: " example.com , google.com , github.com ",
			want:  []string{"example.com", "google.com", "github.com"},
		},
		{
			name:  "domains with empty entries",
			input: "example.com,,google.com,",
			want:  []string{"example.com", "google.com"},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only commas",
			input:   ",,,",
			wantErr: true,
		},
		{
			name:    "invalid domain format with spaces",
			input:   "example.com,host with spaces",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDomainsFromString(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseDomainsFromString() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDomainsFromString() unexpected error: %v", err)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ParseDomainsFromString() length = %d, want %d", len(got), len(tt.want))
				return
			}

			for i, domain := range got {
				if domain != tt.want[i] {
					t.Errorf("ParseDomainsFromString()[%d] = %v, want %v", i, domain, tt.want[i])
				}
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "valid hostname",
			input: "example.com",
		},
		{
			name:  "valid hostname with port",
			input: "example.com:443",
		},
		{
			name:  "hostname with subdomain",
			input: "www.example.com",
		},
		{
			name:  "IPv6 address with brackets",
			input: "[::1]:8080",
		},
		{
			name:  "IPv6 address with brackets no port",
			input: "[::1]",
		},
		{
			name:  "IPv6 address without brackets",
			input: "2001:db8::1",
		},
		{
			name:  "valid port range",
			input: "example.com:65535",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "only spaces",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "hostname with spaces",
			input:   "exam ple.com",
			wantErr: true,
		},
		{
			name:    "empty hostname with port",
			input:   ":443",
			wantErr: true,
		},
		{
			name:    "invalid port - non-numeric",
			input:   "example.com:abc",
			wantErr: true,
		},
		{
			name:    "invalid port - too low",
			input:   "example.com:0",
			wantErr: true,
		},
		{
			name:    "invalid port - too high",
			input:   "example.com:65536",
			wantErr: true,
		},
		{
			name:    "IPv6 missing closing bracket",
			input:   "[::1:8080",
			wantErr: true,
		},
		{
			name:    "IPv6 empty address",
			input:   "[]",
			wantErr: true,
		},
		{
			name:    "IPv6 invalid format after bracket",
			input:   "[::1]invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHost(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateHost() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("validateHost() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "ssl-cert-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test valid config file
	validConfigPath := filepath.Join(tempDir, "valid.yaml")
	validConfigContent := `hosts:
  - example.com
  - google.com:443
  - github.com
`
	if err := os.WriteFile(validConfigPath, []byte(validConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	config, err := LoadConfig(validConfigPath)
	if err != nil {
		t.Errorf("LoadConfig() unexpected error: %v", err)
	}
	if config == nil {
		t.Fatal("LoadConfig() returned nil config")
	}
	if len(config.Hosts) != 3 {
		t.Errorf("LoadConfig() hosts count = %d, want 3", len(config.Hosts))
	}

	// Test empty config file
	emptyConfigPath := filepath.Join(tempDir, "empty.yaml")
	if err := os.WriteFile(emptyConfigPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	_, err = LoadConfig(emptyConfigPath)
	if err == nil {
		t.Error("LoadConfig() expected error for empty file but got none")
	}

	// Test invalid YAML
	invalidConfigPath := filepath.Join(tempDir, "invalid.yaml")
	invalidConfigContent := `hosts:
  - example.com
  invalid yaml content
`
	if err := os.WriteFile(invalidConfigPath, []byte(invalidConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = LoadConfig(invalidConfigPath)
	if err == nil {
		t.Error("LoadConfig() expected error for invalid YAML but got none")
	}

	// Test non-existent file
	_, err = LoadConfig("/non/existent/file.yaml")
	if err == nil {
		t.Error("LoadConfig() expected error for non-existent file but got none")
	}

	// Test empty path
	_, err = LoadConfig("")
	if err == nil {
		t.Error("LoadConfig() expected error for empty path but got none")
	}
}

func TestAppConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  AppConfig
		wantErr bool
	}{
		{
			name: "valid config with config file",
			config: AppConfig{
				ConfigFile:   "config.yaml",
				Timeout:      5,
				OutputFormat: "table",
			},
		},
		{
			name: "valid config with domains",
			config: AppConfig{
				Domains:      "example.com,google.com",
				Timeout:      10,
				OutputFormat: "json",
			},
		},
		{
			name: "valid config with empty output format",
			config: AppConfig{
				Domains: "example.com",
				Timeout: 5,
			},
		},
		{
			name: "no config or domains",
			config: AppConfig{
				Timeout: 5,
			},
			wantErr: true,
		},
		{
			name: "both config and domains",
			config: AppConfig{
				ConfigFile: "config.yaml",
				Domains:    "example.com",
				Timeout:    5,
			},
			wantErr: true,
		},
		{
			name: "invalid timeout",
			config: AppConfig{
				Domains: "example.com",
				Timeout: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid output format",
			config: AppConfig{
				Domains:      "example.com",
				Timeout:      5,
				OutputFormat: "xml",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("AppConfig.Validate() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("AppConfig.Validate() unexpected error: %v", err)
				}
			}
		})
	}
}
