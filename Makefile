# Copyright © 2017. TIBCO Software Inc.
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

HAS_BINDATA := $(shell go-bindata -version 2>/dev/null)

GITBRANCH:=$(shell git rev-parse --abbrev-ref --symbolic-full-name @{u})
MASHLINGLOCALGITREV=`git rev-parse HEAD`
MASHLINGMASTERGITREV=`git rev-parse ${GITBRANCH}`

LDFLAGS= -ldflags "-X ./cli/app.MashlingMasterGitRev=${MASHLINGMASTERGITREV} -X ./cli/app.MashlingLocalGitRev=${MASHLINGLOCALGITREV}  -X ./cli/app.GitBranch=${GITBRANCH}"


.PHONY: all
all: assets install

install:
	rm -f ${GOPATH}/bin/mashling
	go build ${LDFLAGS} ./...
	go install ${LDFLAGS} ./...
	mashling version

assets:
ifndef HAS_BINDATA
	go get github.com/jteeuwen/go-bindata/...
endif
	cd cli && go-bindata -o assets/assets.go -pkg assets \
	assets/banner.txt \
	assets/defGopkg.lock \
	assets/defGopkg.toml \
	schema/mashling_schema-0.2.json
