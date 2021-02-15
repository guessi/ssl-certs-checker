FROM golang:1.15-alpine3.13 AS builder
LABEL maintainer="guessi <guessi@gmail.com>"
RUN apk add --no-cache git ca-certificates
WORKDIR ${GOPATH}/src/github.com/guessi/ssl-certs-checker
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/ssl-certs-checker

FROM scratch
LABEL maintainer="guessi <guessi@gmail.com>"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/ssl-certs-checker /opt/
WORKDIR /opt/
ENTRYPOINT ["/opt/ssl-certs-checker"]
