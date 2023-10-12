FROM public.ecr.aws/docker/library/golang:1.21-alpine3.18 AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR ${GOPATH}/src/github.com/guessi/ssl-certs-checker
COPY *.go go.mod go.sum ./
RUN GOPROXY=direct GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/ssl-certs-checker

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/ssl-certs-checker /opt/
WORKDIR /opt/
ENTRYPOINT ["/opt/ssl-certs-checker"]
