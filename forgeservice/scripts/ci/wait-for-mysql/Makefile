.PHONY: build

GO := CGO_ENABLED=0 go

LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

build:
	@echo ">> building for local environment and for linux"
	$(GO) build --installsuffix cgo -v -ldflags '-w $(LDFLAGS)' wait-for-mysql.go
	GOOS=linux GOARCH=amd64 $(GO) build --installsuffix cgo -v -ldflags '-w $(LDFLAGS)' -o wait-for-mysql.linux wait-for-mysql.go
