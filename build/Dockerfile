# Stage 1
ARG baseImage="golang:1.15-alpine"

FROM ${baseImage} AS build

ENV CGO_ENABLED=0 GOOS=linux

COPY . /app

WORKDIR /app

RUN go build -o spa-server ./cmd/spa-server



# Stage 2 - Build minimised image
FROM node:12-alpine

RUN apk add --no-cache curl=7.67.0-r1

COPY --from=build /app/spa-server /app/bin/spa-server

COPY certs /etc/spa-server/certs

COPY configs/config.default.yaml /etc/spa-server/config.yaml

COPY web /var/www/html

HEALTHCHECK --interval=15s --timeout=5s CMD (curl -k -X GET -Is http://localhost:8079 | grep "204") || exit 1

ENTRYPOINT ["/app/bin/spa-server"]
