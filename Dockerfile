FROM mhart/alpine-node:6.4.0

MAINTAINER TIBCO Software Inc.

RUN apk add --no-cache ca-certificates

ENV GOLANG_VERSION 1.7
ENV GOLANG_SRC_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_SRC_SHA256 72680c16ba0891fcf2ccf46d0f809e4ecf47bbf889f5d884ccb54c5e9a17e1c0

# https://golang.org/issue/14851
COPY go/1.7/no-pic.patch /

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
	\
	&& export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
	\
	&& wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
	&& echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& cd /usr/local/go/src \
	&& patch -p2 -i /no-pic.patch \
	&& ./make.bash \
	\
	&& rm -rf /*.patch \
	&& apk del .build-deps

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH" \
  && apk add --no-cache bash git \
  && echo "Installing GB" \
  && go get -u github.com/constabulary/gb/...
  && go get -u github.com/TIBCOSoftware/flogo-cli/...
  && go get -u github.com/TIBCOSoftware/flogo-lib/...



COPY go/1.7/go-wrapper /usr/local/bin/


VOLUME ["/data"]

EXPOSE 4404
echo "Port 4404 is up and running"