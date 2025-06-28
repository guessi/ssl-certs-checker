package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/guessi/ssl-certs-checker/pkg/cert"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// NewFormatter creates a new output formatter
func New() *Formatter {
	return &Formatter{}
}

// Format formats the certificate results according to the specified format
func (f *Formatter) Format(result *cert.Result, format string) error {
	switch format {
	case "json":
		return f.formatJSON(result)
	case "yaml":
		return f.formatYAML(result)
	case "table", "":
		return f.formatTable(result)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// formatJSON outputs the results in JSON format
func (f *Formatter) formatJSON(result *cert.Result) error {
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	fmt.Println(string(jsonOutput))
	return nil
}

// formatYAML outputs the results in YAML format
func (f *Formatter) formatYAML(result *cert.Result) error {
	yamlOutput, err := yaml.Marshal(result)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	fmt.Println(string(yamlOutput))
	return nil
}

// formatTable outputs the results in table format
func (f *Formatter) formatTable(result *cert.Result) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Host",
		"Common Name",
		"DNS Names",
		"Not Before",
		"Not After",
		"PublicKeyAlgorithm",
		"Issuer",
	})

	for _, certInfo := range result.Certificates {
		dnsNames := ""
		if len(certInfo.DNSNames) > 0 {
			dnsNames = strings.Join(certInfo.DNSNames, "\n")
		}

		t.AppendRows([]table.Row{{
			certInfo.Host,
			certInfo.CommonName,
			dnsNames,
			certInfo.NotBefore,
			certInfo.NotAfter,
			certInfo.PublicKeyAlgorithm,
			certInfo.Issuer,
		}})
	}

	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "\nErrors encountered:\n")
		for _, errInfo := range result.Errors {
			fmt.Fprintf(os.Stderr, "  %s: %s\n", errInfo.Host, errInfo.Error)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}
