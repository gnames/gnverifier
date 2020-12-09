VERSION = $(shell git describe --tags)
VER = $(shell git describe --tags --abbrev=0)
DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S%Z')
FLAG_MODULE = GO111MODULE=on
FLAGS_SHARED = $(FLAG_MODULE) CGO_ENABLED=0 GOARCH=amd64
FLAGS_LD=-ldflags "-X github.com/gnames/gnverify.Build=${DATE} \
                  -X github.com/gnames/gnverify.Version=${VERSION}"
GOCMD=go
GOINSTALL=$(GOCMD) install $(FLAGS_LD)
GOBUILD=$(GOCMD) build $(FLAGS_LD)
GOCLEAN=$(GOCMD) clean
GOGENERATE=$(GOCMD) generate
GOGET = $(GOCMD) get

all: install

test: deps install
	$(FLAG_MODULE) go test ./...

deps:
	$(GOCMD) mod download; \
	$(GOGENERATE)

build:
	$(GOGENERATE)
	cd gnverify; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) $(GOBUILD);

release:
	cd gnverify; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=linux $(GOBUILD); \
	tar zcvf /tmp/gnverify-${VER}-linux.tar.gz gnverify; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=darwin $(GOBUILD); \
	tar zcvf /tmp/gnverify-${VER}-mac.tar.gz gnverify; \
	$(GOCLEAN); \
	$(FLAGS_SHARED) GOOS=windows $(GOBUILD); \
	zip -9 /tmp/gnverify-${VER}-win-64.zip gnverify.exe; \
	$(GOCLEAN);

install:
	$(GOGENERATE)
	cd gnverify; \
	$(FLAGS_SHARED) $(GOINSTALL);

