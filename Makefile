build_macos:
	GOOS=darwin GOARCH=$(ARCH) go build -o build/cryptctl-darwin-$(ARCH) -ldflags "-X github.com/shubhindia/cryptctl/commands.Version=$(git describe --tags)" -a .

build_linux: