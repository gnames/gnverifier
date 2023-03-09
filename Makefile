PROJ_NAME = gnverifier

VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')

NO_C = CGO_ENABLED=0
FLAGS_SHARED = $(NO_C) GOARCH=amd64
FLAGS_LD=-ldflags "-w -s \
                  -X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE) \
                  -X github.com/gnames/$(PROJ_NAME)/pkg.Version=$(VERSION)"
FLAGS_REL = -trimpath -ldflags "-s -w \
						-X github.com/gnames/$(PROJ_NAME)/pkg.Build=$(DATE)"

GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GORELEASE = $(GOCMD) build $(FLAGS_REL)
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get

all: install

test: deps install
	@echo Run tests
	$(GOCMD) test -shuffle=on -count=1 -race -coverprofile=coverage.txt ./...

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

deps:
	@echo Download go.mod dependencies
	$(GOCMD) mod download; \
	$(GOGENERATE)

build:
	@echo Building
	$(GOGENERATE)
	$(GOCLEAN); \
	$(NO_C) $(GOBUILD);

buildrel:
	@echo Building release binary
	$(GOGENERATE)
	$(GOCLEAN); \
	$(NO_C) $(GORELEASE);

install:
	@echo Build and install locally
	$(GOGENERATE)
	$(NO_C) $(GOINSTALL);

release: dockerhub
	@echo Building releases for Linux, Mac, Windows
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/$(PROJ_NAME)-${VER}-linux.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=darwin $(GOBUILD); \
	tar zcvf /tmp/$(PROJ_NAME)-${VER}-mac.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(NO_C) GOOS=darwin GOARCH=arm64 $(GOBUILD); \
	tar zcvf /tmp/$(PROJ_NAME)-${VER}-mac-arm64.tar.gz $(PROJ_NAME); \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=windows $(GOBUILD); \
	zip -9 /tmp/$(PROJ_NAME)-${VER}-win-64.zip $(PROJ_NAME).exe; \
	$(GOCLEAN);

docker: build
	@echo Build Docker images
	docker buildx build -t gnames/$(PROJ_NAME):latest -t gnames/$(PROJ_NAME):$(VERSION) .; \
	$(GOCLEAN);

dockerhub: docker
	@echo Push Docker images to DockerHub
	docker push gnames/$(PROJ_NAME); \
	docker push gnames/$(PROJ_NAME):$(VERSION)

