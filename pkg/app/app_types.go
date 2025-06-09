package app

import (
	"github.com/guessi/ssl-certs-checker/pkg/cert"
	"github.com/guessi/ssl-certs-checker/pkg/output"
)

type App struct {
	checker   *cert.Checker
	formatter *output.Formatter
}
