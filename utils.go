package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Hosts []string `yaml:"hosts"`
}

func readConfig(config string) Config {
	c := Config{}

	y, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Printf("fatal: %s\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(y, &c)
	if err != nil {
		fmt.Printf("fatal: %s\n", err)
		os.Exit(1)
	}
	return c
}

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

func getCells(t table.Writer, host string, port int) {
	certs, err := getPeerCertificates(host, port)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return // skip if target host invalid
	}

	for _, c := range certs {
		if c.IsCA {
			continue
		}
		t.AppendRows([]table.Row{{
			host + ":" + strconv.Itoa(port),
			(*c).Subject.CommonName,
			strings.Join((*c).DNSNames, "\n"),
			(*c).NotBefore,
			(*c).NotAfter,
			(*c).Issuer.CommonName,
		}})
	}
}

func prettyPrintCertsInfo(config string) {
	rc := readConfig(config)
	if len(rc.Hosts) <= 0 {
		fmt.Printf("key not found, or empty input\n")
		return
	}

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

	for _, target := range rc.Hosts {
		p := defaultPort
		ts := strings.Split(target, ":")
		if len(ts) == 2 {
			tp, err := strconv.Atoi(ts[1])
			if err != nil {
				fmt.Errorf("err: invalid port [%s], assume target port is 443\n", target)
			} else {
				p = tp
			}
		}

		getCells(t, ts[0], p)
	}

	t.Style().Format.Header = text.FormatDefault
	t.Render()
}
