.PHONY: all build test lint fmt vet clean md5sum

all: build test

GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')

LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./... | grep -v /vendor/)

build:
	@echo ">> building"
	$(GO) build --installsuffix cgo -v -ldflags '$(LDFLAGS)' cmd/service/main.go

build-linux:
	@echo ">> building for linux"
	GOOS=linux $(GO) build --installsuffix cgo -v -ldflags '$(LDFLAGS)' cmd/service/main.go

test:
	@echo ">> running tests"
	@for PKG in $(PACKAGES); do $(GO) test -short -count=1 -cover $$PKG || exit 1; done;

test-resilience:
	@echo ">> running resilience tests"
	@for PKG in $(PACKAGES); do $(GO) test --tags=integration -count=1 -cover $$PKG || exit 1; done;

lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	@for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

fmt:
	$(GO) fmt $(PACKAGES)

vet:
	$(GO) vet $(PACKAGES)

clean:
	$(GO) clean -i ./...
	rm main

