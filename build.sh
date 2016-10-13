#!/bin/sh -ex
ln -s /workspace /go/src/cb
cd /go/src/cb
go get -v
go build -v -o cb
