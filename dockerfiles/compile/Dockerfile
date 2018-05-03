# mashling/mashling-compile
# version 0.4.0
FROM golang:1.10-alpine
MAINTAINER Jeffrey Bozek, jbozek@tibco.com

RUN apk add --update make upx bash git
ENV GOPATH="/mashling"
ENV PATH="${PATH}:/mashling/bin"
WORKDIR /mashling/src/github.com/TIBCOSoftware/mashling
