# SSL certificate checker written in golang


## Setup Guide

    go get -u github.com/guessi/ssl-certs-checker

## Examples
 
install binary to your ${GOPATH}

    go install github.com/ssl-certs-checker

check single target host certificates infomation

    ${GOPATH}/bin/ssl-certs-checker --hosts "www.google.com"

    +----------------+----------------+-------------------------------+-------------------------------+------------+
    | Host           | Common Name    | Not Before                    | Not After                     | Issuer     |
    +----------------+----------------+-------------------------------+-------------------------------+------------+
    | www.google.com | www.google.com | 2020-02-12 11:47:41 +0000 UTC | 2020-05-06 11:47:41 +0000 UTC | GTS CA 1O1 |
    +----------------+----------------+-------------------------------+-------------------------------+------------+

check multiple target hosts' certificates at once

    ${GOPATH}/bin/ssl-certs-checker --hosts "www.google.com,www.azure.com,www.amazon.com"

    +----------------+----------------+-------------------------------+-------------------------------+-----------------------+
    | Host           | Common Name    | Not Before                    | Not After                     | Issuer                |
    +----------------+----------------+-------------------------------+-------------------------------+-----------------------+
    | www.google.com | www.google.com | 2020-02-12 11:47:41 +0000 UTC | 2020-05-06 11:47:41 +0000 UTC | GTS CA 1O1            |
    | www.azure.com  | *.azure.com    | 2019-12-17 19:51:44 +0000 UTC | 2020-12-17 19:51:44 +0000 UTC | Microsoft IT TLS CA 4 |
    | www.amazon.com | www.amazon.com | 2019-09-18 00:00:00 +0000 UTC | 2020-08-23 12:00:00 +0000 UTC | DigiCert Global CA G2 |
    +----------------+----------------+-------------------------------+-------------------------------+-----------------------+

run with docker

    docker build -t ssl-certs-checker .

    docker run --rm -it ssl-certs-checker --hosts "www.google.com"

# License

[MIT LICENSE](LICENSE)
