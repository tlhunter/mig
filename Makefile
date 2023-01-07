MIG_VERSION = $(shell cat version.txt)
BUILD_TIME = $(shell date +"%Y-%m-%dT%H:%M:%S%z")

# Strip debug info
GO_FLAGS += "-ldflags=-s -w -X 'github.com/tlhunter/mig/commands.Version=$(MIG_VERSION)' -X 'github.com/tlhunter/mig/commands.BuildTime=$(BUILD_TIME)'"

build:
	go build $(GO_FLAGS) -o mig

tiny:
	go build $(GO_FLAGS) -o mig
	upx mig

multi:
	echo "Linux amd64 (x86)"
	GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o mig-linux-amd64
	echo "Windows amd64 (x86)"
	GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o mig.exe
	echo "macOS amd64 (Intel)"
	GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) -o mig-macos-amd64
	echo "macOS arm64 (Apple Silicon)"
	GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o mig-macos-arm64

publish:
	go build
	git commit -am "version v$(MIG_VERSION)"
	git tag "v$(MIG_VERSION)"
	git push origin main "v$(MIG_VERSION)"

test:
	go test -v ./...

clean:
	rm mig || true
	rm mig-linux-amd64 || true
	rm mig.exe || true
	rm mig-macos-amd64 || true
	rm mig-macos-arm64 || true
