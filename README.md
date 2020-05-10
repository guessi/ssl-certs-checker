# SSL certificate checker written in golang

[![Docker Stars](https://img.shields.io/docker/stars/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)
[![Docker Pulls](https://img.shields.io/docker/pulls/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)
[![Docker Automated](https://img.shields.io/docker/automated/guessi/ssl-certs-checker.svg)](https://hub.docker.com/r/guessi/ssl-certs-checker/)

## Setup Guide

    go get -u github.com/guessi/ssl-certs-checker

## Examples

run with docker

    docker build -t guessi/ssl-certs-checker .

    docker run --rm -v $(pwd)/hosts.yaml:/opt/hosts.yaml:ro -it guessi/ssl-certs-checker --config hosts.yaml

install binary to your ${GOPATH} and run locally

    go install github.com/ssl-certs-checker

    ${GOPATH}/bin/ssl-certs-checker --config hosts.yaml

    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+--------------------+-----------------------+
    | Host               | Common Name    | DNS Names                             | Not Before                    | Not After                     | PublicKeyAlgorithm | Issuer                |
    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+--------------------+-----------------------+
    | www.google.com:443 | www.google.com | www.google.com                        | 2020-04-15 20:25:31 +0000 UTC | 2020-07-08 20:25:31 +0000 UTC | ECDSA              | GTS CA 1O1            |
    | www.azure.com:443  | *.azure.com    | *.azure.com                           | 2019-12-17 19:51:44 +0000 UTC | 2020-12-17 19:51:44 +0000 UTC | RSA                | Microsoft IT TLS CA 4 |
    | www.amazon.com:443 | www.amazon.com | amazon.com                            | 2020-01-23 00:00:00 +0000 UTC | 2020-12-31 12:00:00 +0000 UTC | RSA                | DigiCert Global CA G2 |
    |                    |                | amzn.com                              |                               |                               |                    |                       |
    |                    |                | buybox.amazon.com                     |                               |                               |                    |                       |
    |                    |                | corporate.amazon.com                  |                               |                               |                    |                       |
    |                    |                | home.amazon.com                       |                               |                               |                    |                       |
    |                    |                | iphone.amazon.com                     |                               |                               |                    |                       |
    |                    |                | konrad-test.amazon.com                |                               |                               |                    |                       |
    |                    |                | mp3recs.amazon.com                    |                               |                               |                    |                       |
    |                    |                | p-nt-www-amazon-com-kalias.amazon.com |                               |                               |                    |                       |
    |                    |                | p-y3-www-amazon-com-kalias.amazon.com |                               |                               |                    |                       |
    |                    |                | p-yo-www-amazon-com-kalias.amazon.com |                               |                               |                    |                       |
    |                    |                | static.amazon.com                     |                               |                               |                    |                       |
    |                    |                | test-www.amazon.com                   |                               |                               |                    |                       |
    |                    |                | uedata.amazon.com                     |                               |                               |                    |                       |
    |                    |                | us.amazon.com                         |                               |                               |                    |                       |
    |                    |                | www.amazon.com                        |                               |                               |                    |                       |
    |                    |                | www.amzn.com                          |                               |                               |                    |                       |
    |                    |                | www.cdn.amazon.com                    |                               |                               |                    |                       |
    |                    |                | www.m.amazon.com                      |                               |                               |                    |                       |
    |                    |                | yellowpages.amazon.com                |                               |                               |                    |                       |
    |                    |                | yp.amazon.com                         |                               |                               |                    |                       |
    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+--------------------+-----------------------+

# License

[MIT LICENSE](LICENSE)
