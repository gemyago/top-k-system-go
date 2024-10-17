baseBuildImage=alpine3.20
baseRuntimeImage=alpine:3.20
goVersion=$(shell grep "^go " ../../go.mod | awk '{print $$2}')
appName=$(shell sed -n 's/^module .*\/\([^/]*\)$$/\1/p' ../../go.mod)
localRegistry=localhost:6000