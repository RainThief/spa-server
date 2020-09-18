# Stage 1
ARG baseImage="golang:1.15"
FROM ${baseImage} AS build
COPY . /app
WORKDIR /app
RUN go build -o spa-server ./cmd/spa-server
ENTRYPOINT /app/spa-server

# Stage 2 - Build minimised image
FROM debian:buster-slim
RUN apt-get update && apt-get -y install --no-install-recommends curl
RUN apt-get autoremove -y && apt-get clean -y
# COPY --from=build /app/go.mod /app/go.mod
COPY --from=build /app/spa-server /app/bin/spa-server
COPY certs /etc/spa-server/certs
COPY configs/config.default.yaml /etc/spa-server/config.yaml
COPY web /var/www/html

HEALTHCHECK --interval=15s --timeout=5s CMD (curl -k -X GET -Is http://localhost:8079 | grep "204") || exit 1

ENTRYPOINT ["/app/bin/spa-server"]
