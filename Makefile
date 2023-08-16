build_macos:
	GOOS=darwin GOARCH=$(ARCH) go build -o build/cryptctl-darwin-$(ARCH) -a .

build_linux: