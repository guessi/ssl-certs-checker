package main

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func getPeerCertificates(h string, port int) ([]*x509.Certificate, error) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{
			Timeout: dialerTimeout,
		},
		protocol,
		h+":"+strconv.Itoa(port),
		&tls.Config{
			ServerName: h,
		})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return nil, err
	}
	return conn.ConnectionState().PeerCertificates, nil
}

func getCells(t table.Writer, h string) {
	certs, err := getPeerCertificates(h, 443)
	if err != nil {
		return // skip if target host invalid
	}

	for _, c := range certs {
		if c.IsCA {
			continue
		}
		t.AppendRows([]table.Row{
			{h, (*c).Subject.CommonName, strings.Join((*c).DNSNames, "\n"), (*c).NotBefore, (*c).NotAfter, (*c).Issuer.CommonName},
		})
	}
}

func prettyPrintCertsInfo(h string) {
	targets := strings.Split(h, ",")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"Host",
		"Common Name",
		"DNS Names",
		"Not Before",
		"Not After",
		"Issuer",
	})

	for _, target := range targets {
		getCells(t, target)
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()
}
