export GOPATH ?= $(HOME)/go
GOPATH_CB = $(GOPATH)/src/github.com/hiroshi/cb
SRCS = ../main.go

all: json-config yaml-config

json-config: node-docker-example.tar.gz
	cd $(GOPATH_CB)/examples && go run $(SRCS) $^ --config cb-demo.json

yaml-config: node-docker-example.tar.gz
	cd $(GOPATH_CB)/examples && go run $(SRCS) $^ --config cb-demo.yml

node-docker-example.tar.gz:
	cd examples && curl -O https://storage.googleapis.com/container-builder-examples/$@
