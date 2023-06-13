build_macos:
	GOOS=darwin GOARCH=$(ARCH) go build -o build/cryptctl-darwin-$(ARCH) -ldflags "-X github.com/shubhindia/cryptctl/common.versionString=0.0.2" -a .

build_linux: