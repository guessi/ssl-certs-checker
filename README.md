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

    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+-----------------------+
    | Host               | Common Name    | DNS Names                             | Not Before                    | Not After                     | Issuer                |
    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+-----------------------+
    | www.google.com:443 | www.google.com | www.google.com                        | 2020-02-12 11:47:41 +0000 UTC | 2020-05-06 11:47:41 +0000 UTC | GTS CA 1O1            |
    | www.azure.com:443  | *.azure.com    | *.azure.com                           | 2019-12-17 19:51:44 +0000 UTC | 2020-12-17 19:51:44 +0000 UTC | Microsoft IT TLS CA 4 |
    | www.amazon.com:443 | www.amazon.com | amazon.com                            | 2019-09-18 00:00:00 +0000 UTC | 2020-08-23 12:00:00 +0000 UTC | DigiCert Global CA G2 |
    |                    |                | amzn.com                              |                               |                               |                       |
    |                    |                | uedata.amazon.com                     |                               |                               |                       |
    |                    |                | us.amazon.com                         |                               |                               |                       |
    |                    |                | www.amazon.com                        |                               |                               |                       |
    |                    |                | www.amzn.com                          |                               |                               |                       |
    |                    |                | corporate.amazon.com                  |                               |                               |                       |
    |                    |                | buybox.amazon.com                     |                               |                               |                       |
    |                    |                | iphone.amazon.com                     |                               |                               |                       |
    |                    |                | yp.amazon.com                         |                               |                               |                       |
    |                    |                | home.amazon.com                       |                               |                               |                       |
    |                    |                | origin-www.amazon.com                 |                               |                               |                       |
    |                    |                | buckeye-retail-website.amazon.com     |                               |                               |                       |
    |                    |                | huddles.amazon.com                    |                               |                               |                       |
    |                    |                | p-nt-www-amazon-com-kalias.amazon.com |                               |                               |                       |
    |                    |                | p-yo-www-amazon-com-kalias.amazon.com |                               |                               |                       |
    |                    |                | p-y3-www-amazon-com-kalias.amazon.com |                               |                               |                       |
    +--------------------+----------------+---------------------------------------+-------------------------------+-------------------------------+-----------------------+

# License

[MIT LICENSE](LICENSE)
