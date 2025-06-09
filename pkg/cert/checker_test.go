package cert

import (
	"context"
	"testing"
	"time"
)

func TestParseHost(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantHostname string
		wantPort     int
		wantErr      bool
	}{
		{
			name:         "hostname only",
			input:        "example.com",
			wantHostname: "example.com",
			wantPort:     443,
			wantErr:      false,
		},
		{
			name:         "hostname with port",
			input:        "example.com:8443",
			wantHostname: "example.com",
			wantPort:     8443,
			wantErr:      false,
		},
		{
			name:         "hostname with default port",
			input:        "example.com:443",
			wantHostname: "example.com",
			wantPort:     443,
			wantErr:      false,
		},
		{
			name:         "hostname with colon but no port",
			input:        "example.com:",
			wantHostname: "example.com",
			wantPort:     443,
			wantErr:      false,
		},
		{
			name:         "IPv6 address with brackets",
			input:        "[::1]:8080",
			wantHostname: "::1",
			wantPort:     8080,
			wantErr:      false,
		},
		{
			name:         "IPv6 address with brackets no port",
			input:        "[::1]",
			wantHostname: "::1",
			wantPort:     443,
			wantErr:      false,
		},
		{
			name:         "IPv6 address with brackets empty port",
			input:        "[::1]:",
			wantHostname: "::1",
			wantPort:     443,
			wantErr:      false,
		},
		{
			name:         "IPv6 address without brackets",
			input:        "2001:db8::1",
			wantHostname: "2001:db8::1",
			wantPort:     443,
			wantErr:      false,
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
			name:    "empty hostname",
			input:   ":443",
			wantErr: true,
		},
		{
			name:    "invalid port",
			input:   "example.com:abc",
			wantErr: true,
		},
		{
			name:    "port out of range - too low",
			input:   "example.com:0",
			wantErr: true,
		},
		{
			name:    "port out of range - too high",
			input:   "example.com:65536",
			wantErr: true,
		},
		{
			name:    "IPv6 invalid bracket format",
			input:   "[::1:missing_bracket",
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
		{
			name:         "hostname with spaces trimmed",
			input:        "  example.com  ",
			wantHostname: "example.com",
			wantPort:     443,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostname, port, err := parseHost(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseHost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseHost() unexpected error: %v", err)
				return
			}

			if hostname != tt.wantHostname {
				t.Errorf("parseHost() hostname = %v, want %v", hostname, tt.wantHostname)
			}

			if port != tt.wantPort {
				t.Errorf("parseHost() port = %v, want %v", port, tt.wantPort)
			}
		})
	}
}

func TestNewChecker(t *testing.T) {
	timeout := 10 * time.Second
	insecure := true

	checker := New(timeout, insecure)

	if checker == nil {
		t.Fatal("NewChecker() returned nil")
	}

	if checker.timeout != timeout {
		t.Errorf("NewChecker() timeout = %v, want %v", checker.timeout, timeout)
	}

	if checker.insecure != insecure {
		t.Errorf("NewChecker() insecure = %v, want %v", checker.insecure, insecure)
	}
}

func TestCheckCertificates_EmptyHosts(t *testing.T) {
	checker := New(5*time.Second, false)
	ctx := context.Background()

	result, err := checker.CheckCertificates(ctx, []string{})

	if err == nil {
		t.Error("CheckCertificates() expected error for empty hosts but got none")
	}

	if result != nil {
		t.Error("CheckCertificates() expected nil result for empty hosts")
	}
}

func TestCheckCertificates_InvalidHosts(t *testing.T) {
	checker := New(5*time.Second, false)
	ctx := context.Background()

	hosts := []string{"", "invalid::host", "host:99999"}
	result, err := checker.CheckCertificates(ctx, hosts)

	if err != nil {
		t.Errorf("CheckCertificates() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("CheckCertificates() returned nil result")
	}

	// Should have errors for all invalid hosts
	if len(result.Errors) != len(hosts) {
		t.Errorf("CheckCertificates() errors count = %d, want %d", len(result.Errors), len(hosts))
	}

	// Should have no successful certificates
	if len(result.Certificates) != 0 {
		t.Errorf("CheckCertificates() certificates count = %d, want 0", len(result.Certificates))
	}
}

func TestCheckCertificates_ContextCancellation(t *testing.T) {
	checker := New(30*time.Second, false) // Long timeout
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately
	cancel()

	hosts := []string{"example.com"}
	result, err := checker.CheckCertificates(ctx, hosts)

	if err == nil {
		t.Error("CheckCertificates() expected context cancellation error but got none")
	}

	if result != nil {
		t.Error("CheckCertificates() expected nil result for cancelled context")
	}
}

// Integration test - only run if we can reach external hosts
func TestCheckCertificates_RealHost(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	checker := New(10*time.Second, false)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a reliable host that should always be available
	hosts := []string{"google.com:443"}
	result, err := checker.CheckCertificates(ctx, hosts)

	if err != nil {
		t.Errorf("CheckCertificates() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("CheckCertificates() returned nil result")
	}

	// Should have at least one certificate or one error
	if len(result.Certificates) == 0 && len(result.Errors) == 0 {
		t.Error("CheckCertificates() returned no certificates and no errors")
	}

	// If we got a certificate, validate its structure
	if len(result.Certificates) > 0 {
		cert := result.Certificates[0]
		if cert.Host == "" {
			t.Error("Certificate host is empty")
		}
		if cert.NotBefore.IsZero() {
			t.Error("Certificate NotBefore is zero")
		}
		if cert.NotAfter.IsZero() {
			t.Error("Certificate NotAfter is zero")
		}
	}
}
func TestFormatAddress(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		port     int
		want     string
	}{
		{
			name:     "regular hostname",
			hostname: "example.com",
			port:     443,
			want:     "example.com:443",
		},
		{
			name:     "IPv4 address",
			hostname: "192.168.1.1",
			port:     8080,
			want:     "192.168.1.1:8080",
		},
		{
			name:     "IPv6 address",
			hostname: "2001:db8::1",
			port:     443,
			want:     "[2001:db8::1]:443",
		},
		{
			name:     "IPv6 localhost",
			hostname: "::1",
			port:     8443,
			want:     "[::1]:8443",
		},
		{
			name:     "IPv6 already bracketed",
			hostname: "[2001:db8::1]",
			port:     443,
			want:     "[2001:db8::1]:443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAddress(tt.hostname, tt.port)
			if got != tt.want {
				t.Errorf("formatAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
