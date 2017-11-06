# Copyright © 2017. TIBCO Software Inc.
# This file is subject to the license terms contained
# in the license file that is distributed with this file.

HAS_BINDATA := $(shell go-bindata -version 2>/dev/null)

all:
ifndef HAS_BINDATA
	go get github.com/jteeuwen/go-bindata/...
endif
	go-bindata -o lib/model/bindata.go -pkg model lib/model/data/
