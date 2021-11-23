# SPA Server
@todo write this

## How to run
`docker run -it -p 80:80 -p 443:443 ghcr.io/rainthief/spa-server/spa-server:latest`


// Only matches if domain is "www.example.com".
r.Host("www.example.com")
// Matches a dynamic subdomain.
r.Host("{subdomain:[a-z]+}.domain.com")

run_dev.sh


document all optional config options and remove from config yaml

if you edit health check and healthcheck port you need to edit docker file and rebuild

CONFIG DEFAULTS
disableHealthCheck: false
healthCheckPort: 8079


it is an spa server so returns main index page, static file found or directory index (optional)
