FROM public.ecr.aws/docker/library/golang:1.24-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apt update && apt install -y git ca-certificates

WORKDIR /app

COPY *.go go.mod go.sum ./
COPY pkg ./pkg

RUN GOPROXY=direct GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /go/bin/ssl-certs-checker

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/ssl-certs-checker /opt/
WORKDIR /opt/
ENTRYPOINT ["/opt/ssl-certs-checker"]
