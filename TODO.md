

watch file changes
https://stackoverflow.com/questions/19605076/go-auto-recompile-and-reload-server-on-file-change
https://github.com/cespare/reflex


tidy bash scripts

test server.go line 111

Go modules location on co but not in dev cont
Use domain to generate certs
To-do list


bash colours to be identified by context spa server


code coverage


# @todo
# logLevel is the level of logging for this application (DEBUG, INFO, WARN, ERROR)
appLogLevel: "INFO"
httpLog: true
tlsEnable: true
redirectNonTls: true
redirectFromPort: 80


setExpires:
  png: "1 month"
  jpg: "1 month"
httpReadTimeout: 15s
httpWriteTimeout: 15s


writerobottests: true
writeunittests: true
setupci: true



# env overrides config
# config in /internal  and /pkg (main file parsing map string interface)
# env vars overwwrite config, in case of spadir it will overwrite spadirs with only 1 directory
# html dir needs to be web/html https://github.com/golang-standards/project-layout
# config.yaml needs to be in /configs
# dockerfile in /build
# /githooks

gen ssl certs on run_dev
@todo update readme with where certs are


hadolint

shellcheck


cleanup bash

fix tagging

Sigkilll on timeout or does dicker do this
Make list of features and set release point

binary artifact in relases page
