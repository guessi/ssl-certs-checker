package cert

import (
	"time"
)

type CertificateInfo struct {
	Host               string    `json:"host"`
	CommonName         string    `json:"common_name"`
	DNSNames           []string  `json:"dns_names"`
	NotBefore          time.Time `json:"not_before"`
	NotAfter           time.Time `json:"not_after"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm"`
	Issuer             string    `json:"issuer"`
}

type ErrorInfo struct {
	Host  string `json:"host"`
	Error string `json:"error"`
}

type Result struct {
	Certificates []CertificateInfo `json:"certificates"`
	Errors       []ErrorInfo       `json:"errors,omitempty"`
}

type Checker struct {
	timeout  time.Duration
	insecure bool
}
