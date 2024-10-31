# Description: Docker variables for the build process

# base images can be adjusted as needed
baseBuildImage=alpine3.20
runtimeImage=alpine:3.20

# used for local purposes
localRegistry=localhost:6000

# derived from project
goVersion=$(shell grep "^go " ../../go.mod | awk '{print $$2}')
appName=$(shell sed -n 's/^module .*\/\([^/]*\)$$/\1/p' ../../go.mod)