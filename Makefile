.PHONY: build clean all
.SILENT: build clean

ARCH ?= amd64
GO_HOME ?= /go/src/github.com/jinghzhu/KubernetesCRD
ALL_ARCH = amd64 dawrin
GO_VERSION ?= 1.9
GO_IMAGE ?= golang
GO_BIN ?= $(GO_HOME)/bin/$(ARCH)
CMDS = github.com/jinghzhu/KubernetesCRD/cmd/crd

all: build

build:
	for command in $(CMDS) ; do \
		echo "building $$command......."; \
		docker run --rm -u $$(id -u):$$(id -g) -v $$(pwd):$(GO_HOME) \
			-it $(GO_IMAGE):$(GO_VERSION) \
			/bin/sh -c "\
				mkdir -p $(GO_BIN) && \
				GOBIN=$(GO_BIN) go install $$command " && \
		BIN=$$(basename $$command) && echo "Generated bin/$(ARCH)/$$BIN" ; \
	done

clean:
	rm -rf $(pwd)/bin/$(ARCH)
	docker rmi crd
