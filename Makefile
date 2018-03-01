# Copyright © 2017. TIBCO Software Inc.
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

HAS_BINDATA := $(shell go-bindata -version 2>/dev/null)

GITBRANCH:=$(shell git rev-parse --abbrev-ref --symbolic-full-name @{u})
MASHLINGLOCALGITREV=`git rev-parse HEAD`
MASHLINGMASTERGITREV=`git rev-parse ${GITBRANCH}`
#FLOGOGITREV=`git --git-dir=../flogo-lib/.git rev-parse HEAD`

#LDFLAGS= -ldflags "-X github.com/TIBCOSoftware/mashling/cli/app.MashlingMasterGitRev=${MASHLINGMASTERGITREV} -X github.com/TIBCOSoftware/mashling/cli/app.FlogoGitRev=${FLOGOGITREV} -X github.com/TIBCOSoftware/mashling/cli/app.MashlingLocalGitRev=${MASHLINGLOCALGITREV}  -X github.com/TIBCOSoftware/mashling/cli/app.GitBranch=${GITBRANCH}"
LDFLAGS= -ldflags "-X github.com/TIBCOSoftware/mashling/cli/app.MashlingMasterGitRev=${MASHLINGMASTERGITREV} -X github.com/TIBCOSoftware/mashling/cli/app.MashlingLocalGitRev=${MASHLINGLOCALGITREV}  -X github.com/TIBCOSoftware/mashling/cli/app.GitBranch=${GITBRANCH}"


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
