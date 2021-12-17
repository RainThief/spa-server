FROM golang:1.16-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.sum go.sum
COPY go.mod go.mod

RUN --mount=type=cache,target=/var/cache/go go build -o spa-server ./cmd/spa-server



FROM alpine:3.15

RUN apk add --no-cache curl=7.80.0-r0

COPY --from=builder /app/spa-server /app/bin/spa-server

COPY certs /etc/spa-server/certs

COPY configs/config.default.yaml /etc/spa-server/config.yaml

COPY web /var/www/html

# HEALTHCHECK --interval=15s --timeout=5s CMD (curl -k -X GET -Is http://localhost:8000 | grep "204") || exit 1

ENTRYPOINT ["/app/bin/spa-server"]
