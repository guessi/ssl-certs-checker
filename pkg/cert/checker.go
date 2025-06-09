package cert

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultPort    = 443
	Protocol       = "tcp"
	MaxConcurrency = 10
)

// New creates a new certificate checker
func New(timeout time.Duration, insecure bool) *Checker {
	return &Checker{
		timeout:  timeout,
		insecure: insecure,
	}
}

// CheckCertificates checks SSL certificates for multiple hosts concurrently
func (c *Checker) CheckCertificates(ctx context.Context, hosts []string) (*Result, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts provided")
	}

	result := &Result{
		Certificates: make([]CertificateInfo, 0),
		Errors:       make([]ErrorInfo, 0),
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Limit concurrent connections to be respectful to target servers
	semaphore := make(chan struct{}, MaxConcurrency)

	for _, hostStr := range hosts {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		hostname, port, err := parseHost(hostStr)
		if err != nil {
			mutex.Lock()
			result.Errors = append(result.Errors, ErrorInfo{
				Host:  hostStr,
				Error: fmt.Sprintf("invalid host format: %v", err),
			})
			mutex.Unlock()
			continue
		}

		wg.Add(1)
		go func(host string, p int) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			certInfo, err := c.getCertInfoByHost(ctx, host, p)

			mutex.Lock()
			if err != nil {
				result.Errors = append(result.Errors, ErrorInfo{
					Host:  fmt.Sprintf("%s:%d", host, p),
					Error: err.Error(),
				})
			} else if certInfo != nil {
				result.Certificates = append(result.Certificates, *certInfo)
			}
			mutex.Unlock()
		}(hostname, port)
	}

	wg.Wait()
	return result, nil
}

// getCertInfoByHost get SSL certificate info by host
func (c *Checker) getCertInfoByHost(ctx context.Context, hostname string, port int) (*CertificateInfo, error) {
	if hostname == "" {
		return nil, fmt.Errorf("hostname cannot be empty")
	}

	certs, err := c.getPeerCertificates(ctx, hostname, port)
	if err != nil {
		return nil, err
	}

	// Find the first non-CA certificate (leaf certificate)
	for _, cert := range certs {
		if cert == nil || cert.IsCA {
			continue
		}

		return &CertificateInfo{
			Host:               fmt.Sprintf("%s:%d", hostname, port),
			CommonName:         cert.Subject.CommonName,
			DNSNames:           cert.DNSNames,
			NotBefore:          cert.NotBefore,
			NotAfter:           cert.NotAfter,
			PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
			Issuer:             cert.Issuer.CommonName,
		}, nil
	}

	return nil, fmt.Errorf("no valid leaf certificate found")
}

// getPeerCertificates retrieves raw certificates from the server
func (c *Checker) getPeerCertificates(ctx context.Context, hostname string, port int) ([]*x509.Certificate, error) {
	// Create a context with timeout for the entire operation
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	dialer := &net.Dialer{
		Timeout: c.timeout,
	}

	tlsConfig := &tls.Config{
		ServerName:         hostname,
		InsecureSkipVerify: c.insecure,
	}

	address := formatAddress(hostname, port)

	// Use DialContext to respect context cancellation
	conn, err := tls.DialWithDialer(dialer, Protocol, address, tlsConfig)
	if err != nil {
		// Check if the error is due to context cancellation
		select {
		case <-ctxWithTimeout.Done():
			return nil, fmt.Errorf("connection to %s timed out or was cancelled: %w", address, ctxWithTimeout.Err())
		default:
			return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
		}
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			// Log the close error, but don't override the main error
			// In a production environment, you might want to use a proper logger here
		}
	}()

	// Check context cancellation before handshake
	select {
	case <-ctxWithTimeout.Done():
		return nil, ctxWithTimeout.Err()
	default:
	}

	if err := conn.Handshake(); err != nil {
		return nil, fmt.Errorf("TLS handshake failed for %s: %w", address, err)
	}

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no peer certificates found for %s", address)
	}

	return certs, nil
}

// formatAddress formats hostname and port into a proper address string
func formatAddress(hostname string, port int) string {
	// Check if hostname contains colons (potential IPv6)
	if strings.Contains(hostname, ":") && !strings.HasPrefix(hostname, "[") {
		// Likely IPv6 address, wrap in brackets
		return fmt.Sprintf("[%s]:%d", hostname, port)
	}
	return fmt.Sprintf("%s:%d", hostname, port)
}

// parseHost parses a host string into hostname and port
func parseHost(hostStr string) (hostname string, port int, err error) {
	hostStr = strings.TrimSpace(hostStr)
	if hostStr == "" {
		return "", 0, fmt.Errorf("host cannot be empty")
	}

	// Handle IPv6 addresses with brackets [::1]:8080
	if strings.HasPrefix(hostStr, "[") {
		closeBracket := strings.Index(hostStr, "]")
		if closeBracket == -1 {
			return "", 0, fmt.Errorf("invalid IPv6 address format (missing closing bracket): %s", hostStr)
		}

		hostname = hostStr[1:closeBracket]
		if hostname == "" {
			return "", 0, fmt.Errorf("IPv6 address cannot be empty")
		}

		remainder := hostStr[closeBracket+1:]
		if remainder == "" {
			return hostname, DefaultPort, nil
		}

		if !strings.HasPrefix(remainder, ":") {
			return "", 0, fmt.Errorf("invalid format after IPv6 address: %s", hostStr)
		}

		portStr := strings.TrimSpace(remainder[1:])
		if portStr == "" {
			return hostname, DefaultPort, nil
		}

		p, err := strconv.Atoi(portStr)
		if err != nil {
			return "", 0, fmt.Errorf("invalid port number: %s", portStr)
		}
		if p <= 0 || p > 65535 {
			return "", 0, fmt.Errorf("port number out of range (1-65535): %d", p)
		}

		return hostname, p, nil
	}

	// Handle regular hostnames and IPv4 addresses
	parts := strings.Split(hostStr, ":")
	if len(parts) > 2 {
		// Could be IPv6 without brackets, or invalid format
		// Try to determine if it's an IPv6 address
		if strings.Count(hostStr, ":") > 1 {
			// Likely IPv6 address without brackets
			return hostStr, DefaultPort, nil
		}
		return "", 0, fmt.Errorf("invalid host format (too many colons): %s", hostStr)
	}

	hostname = strings.TrimSpace(parts[0])
	if hostname == "" {
		return "", 0, fmt.Errorf("hostname cannot be empty")
	}

	port = DefaultPort
	if len(parts) == 2 {
		portStr := strings.TrimSpace(parts[1])
		if portStr != "" {
			p, err := strconv.Atoi(portStr)
			if err != nil {
				return "", 0, fmt.Errorf("invalid port number: %s", portStr)
			}
			if p <= 0 || p > 65535 {
				return "", 0, fmt.Errorf("port number out of range (1-65535): %d", p)
			}
			port = p
		}
	}

	return hostname, port, nil
}
