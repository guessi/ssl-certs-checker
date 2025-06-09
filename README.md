# SSL Certificate Checker

[![GoDoc](https://godoc.org/github.com/guessi/ssl-certs-checker?status.svg)](https://godoc.org/github.com/guessi/ssl-certs-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/ssl-certs-checker)](https://goreportcard.com/report/github.com/guessi/ssl-certs-checker)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guessi/ssl-certs-checker)](https://github.com/guessi/ssl-certs-checker/blob/master/go.mod)

[![Docker Stars](https://img.shields.io/docker/stars/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)
[![Docker Pulls](https://img.shields.io/docker/pulls/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)

A robust, concurrent SSL certificate checker written in Go.

## Usage

### Command Line Options

```bash
docker run --rm -it guessi/ssl-certs-checker --help
```

## Sample Output

```bash
docker run --rm -it guessi/ssl-certs-checker --domains "github.com"
```

```bash
+----------------+-------------+----------------+-------------------------------+-------------------------------+--------------------+------------------------------------------------+
| Host           | Common Name | DNS Names      | Not Before                    | Not After                     | PublicKeyAlgorithm | Issuer                                         |
+----------------+-------------+----------------+-------------------------------+-------------------------------+--------------------+------------------------------------------------+
| github.com:443 | github.com  | github.com     | 2025-02-05 00:00:00 +0000 UTC | 2026-02-05 23:59:59 +0000 UTC | ECDSA              | Sectigo ECC Domain Validation Secure Server CA |
|                |             | www.github.com |                               |                               |                    |                                                |
+----------------+-------------+----------------+-------------------------------+-------------------------------+--------------------+------------------------------------------------+
```

# License

[MIT LICENSE](LICENSE)
