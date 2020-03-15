FROM golang:1.14-alpine3.11 AS builder
RUN apk add --no-cache git
WORKDIR ${GOPATH}/src/github.com/guessi/ssl-certs-checker
COPY . .
RUN go build -o /go/bin/ssl-certs-checker

FROM alpine:3.11
COPY --from=builder /go/bin/ssl-certs-checker /opt/
WORKDIR /opt/
ENTRYPOINT ["/opt/ssl-certs-checker"]
