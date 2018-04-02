# build variables
LAST_TAG 		= $(shell git describe --abbrev=0 --tags)
PREVIOUS_TAG 	= $(shell git describe --abbrev=0 --tags $(LAST_TAG)^)
BUILD_NUMBER 	= $(shell git rev-list master --count)
HASH 			= $(shell git rev-parse --short HEAD)
DATE 			= $(shell go run tools/build-date.go)

# build flags
BUILD_FLAGS = -ldflags "-s -w \
	-X github.com/blomma/viaduct/option.Version=$(LAST_TAG) \
	-X github.com/blomma/viaduct/option.BuildNumber=$(BUILD_NUMBER) \
	-X github.com/blomma/viaduct/option.CommitHash=$(HASH) \
	-X 'github.com/blomma/viaduct/option.CompileDate=$(DATE)'"

EXECUTABLE = viaduct

UNIX_EXECUTABLES = \
	darwin-amd64-$(EXECUTABLE) \
	linux-amd64-$(EXECUTABLE) \
	linux-arm-7-$(EXECUTABLE)
WIN_EXECUTABLES = \
	windows-amd64-$(EXECUTABLE).exe

COMPRESSED_EXECUTABLES = \
	$(UNIX_EXECUTABLES:%=%.tar.bz2) \
	$(WIN_EXECUTABLES:%.exe=%.zip)
COMPRESSED_EXECUTABLE_TARGETS = $(COMPRESSED_EXECUTABLES:%=bin/%)

all: $(EXECUTABLE)

# arm
bin/linux-arm-5-$(EXECUTABLE):
	GOARM=5 GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) -o "$@"
bin/linux-arm-7-$(EXECUTABLE):
	GOARM=7 GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) -o "$@"

# amd64
bin/darwin-amd64-$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) -o "$@"
bin/linux-amd64-$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) -o "$@"
bin/windows-amd64-$(EXECUTABLE).exe:
	GOARCH=amd64 GOOS=windows go build $(BUILD_FLAGS) -o "$@"

%.tar.bz2: %
	tar -jcvf "$<.tar.bz2" "$<"
%.zip: %.exe
	zip "$@" "$<"

$(EXECUTABLE):
	go build $(BUILD_FLAGS) -o "$@"

release-notes:
	git log --format="- %s" $(PREVIOUS_TAG)..$(LAST_TAG)

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: clean
	$(MAKE) $(COMPRESSED_EXECUTABLE_TARGETS)

install:
	go install $(BUILD_FLAGS)

clean:
	rm $(EXECUTABLE) || true
	rm -rf bin/

.PHONY: clean release install
