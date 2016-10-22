export GOPATH ?= $(HOME)/go
GOPATH_CB = $(GOPATH)/src/github.com/hiroshi/cb
SRCS = ../main.go

OSARCH = $(or $(subst _, ,$(OS_ARCH)),$(subst /, ,,$(lastword $(shell go version))))
GOOS ?= $(firstword $(OSARCH))
GOARCH ?= $(lastword $(OSARCH))
CB = bin/$(GOOS)_$(GOARCH)/cb

OS_ARCHS = darwin_amd64 linux_amd64 windows_amd64

# build
build: $(CB)

$(CB): main.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@

build_all: $(foreach x,$(OS_ARCHS),build-$(x))

build-%:
	make build OS_ARCH=$*

clobber:
	rm -rf bin

# release
release: create_release upload_all
create_release:
	curl https://api.github.com/repos/hiroshi/cb/releases \
	  -H "Authorization: token $$GITHUB_TOKEN" \
	  -d '{"tag_name":"latest"}'

RELEASE_ID = $(shell \
  curl -s https://api.github.com/repos/hiroshi/cb/releases/latest \
    -H "Authorization: token $$GITHUB_TOKEN" \
  | jq .id)
ZIP = bin/cb-latest-$(GOOS)-$(GOARCH).zip
upload: $(ZIP)
	curl https://uploads.github.com/repos/hiroshi/cb/releases/$(RELEASE_ID)/assets?name=$(notdir $(ZIP)) \
	  -H "Authorization: token $$GITHUB_TOKEN" \
	  -H "Accept: application/vnd.github.manifold-preview" \
	  -H "Content-Type: application/zip" \
	  --data-binary @$(ZIP)

$(ZIP):
	zip -j $(ZIP) $(CB)

upload_all: $(foreach x,$(OS_ARCHS),upload-$(x))

upload-%:
	make upload OS_ARCH=$*

# quick install
install: symlink
	rm $(GOPATH)/bin/cb
	cd $(GOPATH_CB) && go install -v

# development
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
