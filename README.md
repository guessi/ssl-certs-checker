# SSL certificate checker written in golang

[![GoDoc](https://godoc.org/github.com/guessi/ssl-certs-checker?status.svg)](https://godoc.org/github.com/guessi/ssl-certs-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/ssl-certs-checker)](https://goreportcard.com/report/github.com/guessi/ssl-certs-checker)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/guessi/ssl-certs-checker)](https://github.com/guessi/ssl-certs-checker/blob/master/go.mod)

[![Docker Stars](https://img.shields.io/docker/stars/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)
[![Docker Pulls](https://img.shields.io/docker/pulls/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)

## Usage

    docker run --rm -v $(pwd)/hosts.yaml:/opt/hosts.yaml:ro -it guessi/ssl-certs-checker --help

    NAME:
       SSL Certificate Checker - check SSL certificates at once

    USAGE:
       ssl-certs-checker [global options] command [command options] [arguments...]

    COMMANDS:
       help, h  Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --config value, -C value   config file
       --timeout value, -t value  dialer timeout in second(s) (default: 5)
       --help, -h                 show help (default: false)


## Sample Output

    docker run --rm -v $(pwd)/hosts.yaml:/opt/hosts.yaml:ro -it guessi/ssl-certs-checker --config hosts.yaml

    +--------------------+----------------+----------------+-------------------------------+-------------------------------+--------------------+------------+
    | Host               | Common Name    | DNS Names      | Not Before                    | Not After                     | PublicKeyAlgorithm | Issuer     |
    +--------------------+----------------+----------------+-------------------------------+-------------------------------+--------------------+------------+
    | www.google.com:443 | www.google.com | www.google.com | 2021-01-19 08:04:07 +0000 UTC | 2021-04-13 08:04:06 +0000 UTC | ECDSA              | GTS CA 1O1 |
    +--------------------+----------------+----------------+-------------------------------+-------------------------------+--------------------+------------+

## Build from Source

    go get -u github.com/guessi/ssl-certs-checker

    cd ${GOPATH}/src/github.com/guessi/ssl-certs-checker

    vim ... # made some changes

    go install github.com/guessi/ssl-certs-checker

    ssl-certs-checker --help

# License

[MIT LICENSE](LICENSE)
