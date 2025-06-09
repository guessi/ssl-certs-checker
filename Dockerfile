FROM public.ecr.aws/docker/library/golang:1.24-alpine3.22 AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache git ca-certificates
WORKDIR ${GOPATH}/src/github.com/guessi/ssl-certs-checker
COPY *.go go.mod go.sum ./
RUN GOPROXY=direct GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /go/bin/ssl-certs-checker

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/ssl-certs-checker /opt/
WORKDIR /opt/
ENTRYPOINT ["/opt/ssl-certs-checker"]
