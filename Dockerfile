# Stage 1
FROM golang:1.14 AS build
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
COPY config.default.yaml /etc/spa-server/config.yaml
COPY web /var/www/html

# @todo needs to get port dynamically
HEALTHCHECK --interval=15s --timeout=5s CMD (curl -k -X GET -Is https://localhost:80 | grep "200") || exit 1
# @todo remove expose
EXPOSE 8080
ENTRYPOINT ["/app/bin/spa-server"]
