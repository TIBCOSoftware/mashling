# Copyright © 2017. TIBCO Software Inc.
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

HAS_BINDATA := $(shell go-bindata -version 2>/dev/null)

VERSION=`git tag | sort -n | tail -1`
MASHLINGLOCALGITREV=`git rev-parse HEAD`
MASHLINGMASTERGITREV=`git rev-parse origin/master`
FLOGOGITREV=`git --git-dir=../flogo-lib/.git rev-parse HEAD`

LDFLAGS= -ldflags "-X github.com/TIBCOSoftware/mashling/cli/app.Version=${VERSION} -X github.com/TIBCOSoftware/mashling/cli/app.MashlingMasterGitRev=${MASHLINGMASTERGITREV} -X github.com/TIBCOSoftware/mashling/cli/app.FlogoGitRev=${FLOGOGITREV} -X github.com/TIBCOSoftware/mashling/cli/app.MashlingLocalGitRev=${MASHLINGLOCALGITREV}"


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
	assets/default_manifest \
	schema/mashling_schema-0.2.json
