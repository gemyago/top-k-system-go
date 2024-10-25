# taken from https://gist.github.com/sighingnow/deee806603ec9274fd47
# below may need to be extended to handle more cases (e.g M1)
UNAME_S := $(shell uname -s)
UNAME_P = $(shell uname -p)
ifeq ($(UNAME_S),Linux)
	os=linux
else
	os=darwin
endif
ifeq ($(UNAME_P),arm)
	arch=arm64
else
	arch=amd64
endif

tar_version := $(shell tar --version)
ifeq ($(findstring GNU tar,$(tar_version)),GNU tar)
	tar_ignore_arg=--ignore-failed-read
	tar_absolute_names=--absolute-names
else
	tar_ignore_arg=
	tar_absolute_names=
endif

bin=bin
tmp=./tmp

# See https://github.com/helm/helm/releases for available versions
helm_version=$(shell cat .helm-version)
helm_archive_name=helm-v$(helm_version)-$(os)-$(arch).tar.gz
helm_archive=$(tmp)/$(helm_archive_name)
helm=$(bin)/helm

# See https://github.com/helmfile/helmfile/releases for available versions
helmfile_version=$(shell cat .helmfile-version)
helmfile_archive_name=helmfile_$(helmfile_version)_$(os)_$(arch).tar.gz
helmfile_archive=$(tmp)/$(helmfile_archive_name)
helmfile_url="https://github.com/helmfile/helmfile/releases/download/v$(helmfile_version)/$(helmfile_archive_name)"
helmfile=$(bin)/helmfile
helmfile_args_file=$(tmp)/helmfile-args.txt

# See https://kubernetes.io/releases/ for available versions
kubectl_version=$(shell cat .kubectl-version)
kubectl_url="https://dl.k8s.io/release/v$(kubectl_version)/bin/$(os)/$(arch)/kubectl"
kubectl=$(bin)/kubectl

# Deploying to local by default
DEPLOY_ENV ?= default

$(tmp):
	mkdir -p $(tmp)

$(bin):
	mkdir -p $(bin)

$(helm_archive): .helm-version $(tmp)
	curl https://get.helm.sh/$(helm_archive_name) -o $(helm_archive)

$(helm): $(helm_archive) $(bin)
	tar -C $(tmp) -xvf $(helm_archive)
	cp $(tmp)/$(os)-$(arch)/helm $(helm)
	touch $(helm)

$(helmfile_archive): .helmfile-version $(tmp)
	curl -L $(helmfile_url) -o $(helmfile_archive)

$(helmfile): $(helmfile_archive) $(bin)
	tar -C $(tmp) -xvf $(helmfile_archive) $(tar_ignore_arg) $(tar_absolute_names)
	cp $(tmp)/helmfile $(helmfile)
	touch $(helmfile)

$(kubectl): .kubectl-version $(bin)
	curl -L $(kubectl_url) -o $(kubectl)
	chmod +x $(kubectl)
	touch $(kubectl)

tools: $(helm) $(helmfile) $(kubectl)
	$(helm) plugin install https://github.com/databus23/helm-diff

.PHONY: $(helmfile_args_file)
$(helmfile_args_file): 
	@echo "-e $(DEPLOY_ENV)" > $@
#	@echo "--debug" >> $@

diff: $(helmfile) $(helmfile_args_file)
	$(helmfile) $(shell cat $(helmfile_args_file)) diff

apply: $(helmfile) $(helmfile_args_file)
	$(helmfile) $(shell cat $(helmfile_args_file)) apply

helmfile: $(helmfile) $(helmfile_args_file)
	@echo $(helmfile) $(shell cat $(helmfile_args_file))
