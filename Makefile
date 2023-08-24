version=$(shell git describe --tags)

build_macos:
	GOOS=darwin GOARCH=$(ARCH) go build -o build/cryptctl-darwin-$(ARCH) -ldflags "-X github.com/shubhindia/cryptctl/commands.Version=${version}" -a .

build_linux:
	GOOS=linux GOARCH=$(ARCH) go build -o build/cryptctl-linux-$(ARCH) -ldflags "-X github.com/shubhindia/cryptctl/commands.Version=${version}" -a .