# git
LAST_TAG := $(shell git describe --abbrev=0 --tags)
PREVIOUS_TAG := $(shell git describe --abbrev=0 --tags $(LAST_TAG)^)

BUILD_NUMBER = $(shell git rev-list master --count)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell go run tools/build-date.go)

# go
BUILD_FLAGS := -ldflags "-s -w -X flag.Version=$(LAST_TAG) -X flag.BuildNumber=$(BUILD_NUMBER) -X flag.CommitHash=$(HASH) -X 'flag.CompileDate=$(DATE)'"

# github-release
USER := blomma

EXECUTABLE := viaduct

UNIX_EXECUTABLES := \
	darwin/amd64/$(EXECUTABLE) \
	linux/amd64/$(EXECUTABLE) \
	linux/arm/7/$(EXECUTABLE)
WIN_EXECUTABLES := \
	windows/amd64/$(EXECUTABLE).exe

COMPRESSED_EXECUTABLES = $(UNIX_EXECUTABLES:%=%.tar.bz2) $(WIN_EXECUTABLES:%.exe=%.zip)
COMPRESSED_EXECUTABLE_TARGETS = $(COMPRESSED_EXECUTABLES:%=bin/%)

UPLOAD_CMD = github-release upload -u $(USER) -r $(EXECUTABLE) -t $(LAST_TAG) -n $(subst /,-,$(FILE)) -f bin/$(FILE)

all: $(EXECUTABLE)

# arm
bin/linux/arm/5/$(EXECUTABLE):
	GOARM=5 GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) -o "$@"
bin/linux/arm/7/$(EXECUTABLE):
	GOARM=7 GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) -o "$@"

# amd64
bin/darwin/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) -o "$@"
bin/linux/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) -o "$@"
bin/windows/amd64/$(EXECUTABLE).exe:
	GOARCH=amd64 GOOS=windows go build $(BUILD_FLAGS) -o "$@"

%.tar.bz2: %
	tar -jcvf "$<.tar.bz2" "$<"
%.zip: %.exe
	zip "$@" "$<"

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: clean
	$(MAKE) $(COMPRESSED_EXECUTABLE_TARGETS)
	git push && git push --tags
	git log --format="- %s" $(PREVIOUS_TAG)..$(LAST_TAG) | \
		github-release release -u $(USER) -r $(EXECUTABLE) \
		-t $(LAST_TAG) -n $(LAST_TAG) -d - || true
	$(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

$(EXECUTABLE):
	go build $(BUILD_FLAGS) -o "$@"

install:
	go install $(BUILD_FLAGS)

clean:
	rm $(EXECUTABLE) || true
	rm -rf bin/

.PHONY: clean release install