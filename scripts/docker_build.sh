#!/usr/bin/env sh
docker run --rm -ti -v $GOPATH:/go -w=/go/src/github.com/opentracing/opentracing-example golang:1.6-alpine go build -v
