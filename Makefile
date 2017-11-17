# Copyright © 2017. TIBCO Software Inc.
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

HAS_BINDATA := $(shell go-bindata -version 2>/dev/null)

VERSION=0.2.0
MASHLINGSCHEMA=0.2
MASHLINGGITTAG=`git rev-parse HEAD`
FLOGOGITTAG=`git --git-dir=../flogo-lib/.git rev-parse HEAD`
LDFLAGS= -ldflags "-X github.com/TIBCOSoftware/mashling/cli/app.Version=${VERSION} -X github.com/TIBCOSoftware/mashling/cli/app.MashlingGitTag=${MASHLINGGITTAG} -X github.com/TIBCOSoftware/mashling/cli/app.ShemaVersion=${MASHLINGSCHEMA} -X github.com/TIBCOSoftware/mashling/cli/app.FlogoGitTag=${FLOGOGITTAG}"

.PHONY: all
all: assets install
	
assets:
ifndef HAS_BINDATA
	go get github.com/jteeuwen/go-bindata/...
endif
	cd cli && go-bindata -o assets/assets.go -pkg assets \
	assets/banner.txt \
	assets/default_manifest \
	schema/mashling_schema-0.2.json

install:
	rm -f ${GOPATH}/bin/mashling
	go build ${LDFLAGS} ./...
	go install ${LDFLAGS} ./...
	mashling version
