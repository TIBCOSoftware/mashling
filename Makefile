IMPORT_PATH := github.com/TIBCOSoftware/mashling
PLATFORMS := darwin/amd64 linux/amd64 linux/arm64 windows/amd64
IGNORED_PACKAGES := /vendor/
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo v0)

GO      = go
GOFMT   = gofmt
DEP     = dep
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

# Support legacy v1 builds
GITBRANCH:=$(shell git rev-parse --abbrev-ref --symbolic-full-name HEAD)
MASHLINGLOCALGITREV=`git rev-parse HEAD`
MASHLINGMASTERGITREV=`git rev-parse ${GITBRANCH}`

.PHONY: build
build: buildgateway buildcli ## Build gateway and cli binaries

.PHONY: all
all: allgateway allcli ## Create assets, run go generate, fmt, vet, and then build then gateway and cli binaries

PRIMARYGOPATH = $(CURDIR)/.GOPATH
SECONDAYGOPATH = $(CURDIR)/.GOPATH/vendor
export GOPATH := $(PRIMARYGOPATH):$(SECONDAYGOPATH)
unexport GOBIN
BIN = $(PRIMARYGOPATH)/bin
export PATH := $(PATH):$(BIN)

# Tools
GOLINT = $(BIN)/golint
$(BIN)/golint: | ; $(info $(M) building golint…)
	$Q go get github.com/golang/lint/golint

GODEP = $(BIN)/dep
$(BIN)/dep: | ; $(info $(M) building dep…)
	$Q go get github.com/golang/dep/cmd/dep

GOMETALINTER = $(BIN)/gometalinter
$(BIN)/gometalinter: | ; $(info $(M) building gometalinter…)
	$Q go get github.com/alecthomas/gometalinter && gometalinter --install

GOBINDATA = go run pkg/assets/bindata.go

.GOPATH/.ok:
	$Q mkdir -p "$(dir .GOPATH/src/$(IMPORT_PATH))"
	$Q ln -s ../../../.. ".GOPATH/src/$(IMPORT_PATH)"
	$Q mkdir -p .GOPATH/vendor/
	$Q ln -s ../../vendor .GOPATH/vendor/src
	$Q mkdir -p bin
	$Q ln -s ../bin .GOPATH/bin
	$Q touch $@

.PHONY: buildgateway
buildgateway: .GOPATH/.ok ; $(info $(M) building gateway executable…) @ ## Build gateway binary
	$Q $(GO) install \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		$(IMPORT_PATH)/cmd/mashling-gateway

.PHONY: allgateway
allgateway: assets generate fmt vet buildgateway ## Satisfy pre-build requirements and then build the gateway binary

.PHONY: buildcli
buildcli: .GOPATH/.ok; $(info $(M) building CLI executable…) @ ## Build CLI binary
	$Q $(GO) install \
		-ldflags '-X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		$(IMPORT_PATH)/cmd/mashling-cli

.PHONY: allcli
allcli: cligenerate cliassets fmt vet buildcli ## Satisfy pre-build requirements and then build the cli binary

.PHONY: buildlegacy
buildlegacy: .GOPATH/.ok; $(info $(M) legacy CLI executable...) @ ## Build legacy CLI binary
	$Q $(GO) install \
		-ldflags "-X github.com/TIBCOSoftware/mashling/cli/app.MashlingMasterGitRev=${MASHLINGMASTERGITREV} \
		-X github.com/TIBCOSoftware/mashling/cli/app.MashlingLocalGitRev=${MASHLINGLOCALGITREV}  \
		-X github.com/TIBCOSoftware/mashling/cli/app.GitBranch=${GITBRANCH}" \
		$(IMPORT_PATH)/cli/cmd/mashling

.PHONY: release
release: .GOPATH/.ok $(PLATFORMS); @ ## Build all executables for release against all targets

$(PLATFORMS): ; $(info $(M) building package executable for $@)
	$Q GOOS=$(os) GOARCH=$(arch) $(GO) build \
		-tags release \
		-ldflags '-s -w -X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
		-o release/mashling-gateway-$(os)-$(arch) $(IMPORT_PATH)/cmd/mashling-gateway && \
		GOOS=$(os) GOARCH=$(arch) $(GO) build \
			-tags release \
			-ldflags '-s -w -X main.Version=$(VERSION) -X main.BuildDate=$(DATE)' \
			-o release/mashling-cli-$(os)-$(arch) $(IMPORT_PATH)/cmd/mashling-cli
		(type -p upx >/dev/null 2>&1 && ( upx release/mashling-gateway-$(os)-$(arch) release/mashling-cli-$(os)-$(arch) ) || echo "UPX not found, skipping compression (please visit https://upx.github.io to install)...")

.PHONY: docker
docker: .GOPATH/.ok linux/amd64 ; $(info $(M) building a docker image containing the mashling-gateway binary) @ ## Build a minimal docker image containing the gateway binary
	$Q type -p docker >/dev/null 2>&1 && docker build . -t mashling-gateway || echo "Docker not found, please visit https://www.docker.com to install for your platform."

.PHONY: setup
setup: clean .GOPATH/.ok gitignoregopath $(GOLINT) $(GODEP) hooks ## Setup the dev environment

gitignoregopath:
	@if ! grep "/.GOPATH" .gitignore > /dev/null 2>&1; then \
	    echo "/.GOPATH" >> .gitignore; \
	    echo "/bin" >> .gitignore; \
	fi

.PHONY: hooks
hooks: .GOPATH/.ok ; $(info $(M) setting up git commit hooks…) @ ## Setup the git commit hooks
	$Q ret=0 && for s in go-fmt go-vet pre-commit; do \
		cp scripts/$$s .git/hooks && chmod +x .git/hooks/$$s || ret=1 ; \
	 done ; exit $$ret

.PHONY: lint
lint: .GOPATH/.ok $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q ret=0 && for pkg in $(allpackages); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	 done ; exit $$ret

.PHONY: metalinter
metalinter: .GOPATH/.ok $(GOMETALINTER) ; $(info $(M) running gometalinter…) @ ## Run gometalinter
	gometalinter ./... --vendor

.PHONY: fmt
fmt: .GOPATH/.ok ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

.PHONY: vet
vet: .GOPATH/.ok ; $(info $(M) running go vet…) @ ## Run go vet on all source files
	$(GO) vet $(allpackages)

.PHONY: generate
generate: .GOPATH/.ok ; $(info $(M) running go generate…) @ ## Run go generate on source
	$(GO) generate ./...

.PHONY: cligenerate
cligenerate: .GOPATH/.ok ; $(info $(M) running CLI go generate…) @ ## Run go generate on CLI source
	$(GO) generate ./internal/app/cli/...

.PHONY: assets
assets: .GOPATH/.ok ; $(info $(M) running asset generation…) @ ## Run asset generation
	$Q { \
	set activities = "" ;\
	set triggers = "" ;\
	set actions = "" ;\
	set assets = "" ;\
	for file in $$(find . -type f -name activity.json -not -path "*.GOPATH*"); do \
	  if [ -e "$${file//.json/.go}" ]; then \
	    echo "Activity Found: $${file//.json/.go}" ;\
	    activities="$$activities $$file" ;\
	  fi; \
	done; \
	for file in $$(find . -type f -name trigger.json -not -path "*.GOPATH*"); do \
	  if [ -e "$${file//.json/.go}" ]; then \
	    echo "Trigger Found: $${file//.json/.go}" ;\
	    triggers="$$triggers $$file" ;\
	  fi ;\
	done ;\
	for file in $$(find . -type f -name action.json -not -path "*.GOPATH*"); do \
	  if [ -e "$${file//.json/.go}" ]; then \
	    echo "Asset Found: $${file//.json/.go} found" ;\
	    actions="$$actions $$file" ;\
	  fi ;\
	done ;\
	for file in $$(find internal/app/gateway/assets/ -not -name "*.go" -type f); do \
		assets="$$assets $$file" ;\
	done ;\
	if [[ ! -z "$$activities" ]]; then \
		$(GOBINDATA) -pkg activities -o internal/app/gateway/flogo/registry/activities/activities.go $$activities ;\
	fi;\
	if [[ ! -z "$$triggers" ]]; then \
		$(GOBINDATA) -pkg triggers -o internal/app/gateway/flogo/registry/triggers/triggers.go $$triggers ;\
	fi;\
	if [[ ! -z "$$actions" ]]; then \
		$(GOBINDATA) -pkg actions -o internal/app/gateway/flogo/registry/actions/actions.go $$actions ;\
	fi;\
	$(GOBINDATA) -prefix internal/app/gateway/assets/ -pkg assets -o internal/app/gateway/assets/assets.go $$assets ;\
	$(GOBINDATA) -prefix cli/ -pkg assets -o cli/assets/assets.go cli/assets/banner.txt cli/assets/defGopkg.lock cli/assets/defGopkg.toml cli/schema/mashling_schema-0.2.json ; \
}

.PHONY: cliassets
cliassets: .GOPATH/.ok ; $(info $(M) running asset generation…) @ ## Run asset generation for CLI
	$Q { \
	set assets = "" ;\
	for file in $$(find internal/app/cli/assets/ -not -name "*.go" -type f); do \
		assets="$$assets $$file" ;\
	done ;\
	$(GOBINDATA) -prefix internal/app/cli/assets/ -pkg assets -o internal/app/cli/assets/assets.go $$assets ;\
}

.PHONY: list
list: .GOPATH/.ok ; $(info $(M) listing internal packages....)	@ ## List packages
	@echo $(allpackages)

# cd into the GOPATH to workaround ./... not following symlinks
_allpackages = $(shell ( cd $(CURDIR)/.GOPATH/src/$(IMPORT_PATH) && \
    GOPATH=$(CURDIR)/.GOPATH go list ./... 2>&1 1>&3 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)) 1>&2 ) 3>&1 | \
    grep -v -e "^$$" $(addprefix -e ,$(IGNORED_PACKAGES)))

# memoize allpackages, so that it's executed only once and only if used
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

.PHONY: dep
dep: .GOPATH/.ok $(GODEP); $(info $(M) verifying and retrieving dependencies…) @ ## Make sure dependencies are vendored
	$Q cd $(PRIMARYGOPATH)/src/$(IMPORT_PATH) && $(DEP) ensure

.PHONY: depadd
depadd: .GOPATH/.ok $(GODEP); $(info $(M) adding dependencies…) @ ## Add new dependencies
	$Q cd $(PRIMARYGOPATH)/src/$(IMPORT_PATH) && $(DEP) ensure -add $(NEWDEPS)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Cleanup everything
	@rm -rf bin release .GOPATH

.PHONY: help
help:
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version:
	@echo $(VERSION)
