export GOPATH ?= $(HOME)/go
GOPATH_CB = $(GOPATH)/src/github.com/hiroshi/cb
SRCS = ../main.go

build:
	go build -o cb

install: symlink
	rm $(GOPATH)/bin/cb
	cd $(GOPATH_CB) && go install -v

goget: symlink
	cd $(GOPATH_CB) && go get -v

symlink: | $(dir $(GOPATH_CB))
	ln -snf $(shell pwd) $(GOPATH_CB)

$(dir $(GOPATH_CB)):
	mkdir -p $@

run-example:
	cd examples && make

# build single binary cb image
cb-build: source.tar.gz
	cd $(GOPATH_CB) && go run main.go $< --config config.yml

source.tar.gz: main.go Dockerfile Dockerfile.build
	tar czvf $@ $^
