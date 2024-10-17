include ../docker.variables.mk

fullLocalBuildRegistry=$(localRegistry)/$(appName)-build
fullLocalOutRegistry=$(localRegistry)/$(appName)

dockerBuildArgs=\
		--build-arg BASE_BUILD_IMAGE=golang:$(goVersion)-$(baseBuildImage) \
		-f Dockerfile.build \
		-t $(fullLocalBuildRegistry):latest \
		../../

.local-build-image:
	docker build $(dockerBuildArgs)
	docker inspect $(fullLocalBuildRegistry):latest > $@-inspect
	cat $@-inspect | jq -r '.[0].RepoTags[] | select(. | startswith("$(fullLocalBuildRegistry)"))' > $@

dockerOutArgs=\
		--build-arg BUILD_IMAGE=$(shell cat .local-build-image) \
		--build-arg RUNTIME_IMAGE=$(runtimeImage) \
		-f Dockerfile.out \
		../../

.local-out-image: .local-build-image
	docker build -t $(fullLocalOutRegistry):latest $(dockerOutArgs)
	docker inspect $(fullLocalOutRegistry):latest > $@-inspect
	cat $@-inspect | jq -r '.[0].RepoTags[] | select(. | startswith("$(fullLocalOutRegistry)"))' > $@

dockerHostArgs=\
		--build-arg RUNTIME_IMAGE=$(runtimeImage) \
		-f Dockerfile.host \
		../../

# NOTE: If building image like this, the platform must be match target execution env while building binaries
# e.g GOOS=linux GOARCH=amd64 make build/bin
.host-local-image:
	docker build -t $(fullLocalOutRegistry):latest $(dockerHostArgs)
	docker inspect $(fullLocalOutRegistry):latest > $@-inspect
	cat $@-inspect | jq -r '.[0].RepoTags[] | select(. | startswith("$(fullLocalOutRegistry)"))' > $@

.PHONY: clean-images
clean-images:
	rm -f .local-*
	rm -f .host-*