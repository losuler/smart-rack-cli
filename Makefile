NAME=srs
VERSION=$(shell git describe --abbrev=0 --tags | cut -c2-)
TAR=tar --gzip --create --verbose
EXTRA_FILES=README.md LICENSE.txt

build:
	go build -o bin/$(NAME) main.go

run:
	go run main.go

build-all:
	GOOS=linux GOOARCH=amd64 go build -v -o build/linux-amd64/$(NAME) main.go
	GOOS=darwin GOARCH=amd64 go build -v -o build/macos-amd64/$(NAME) main.go
	GOOS=windows GOARCH=amd64 go build -v -o build/windows-amd64/$(NAME) main.go

release-all:
	cp $(EXTRA_FILES) build/linux-amd64/
	cp $(EXTRA_FILES) build/macos-amd64/
	cp $(EXTRA_FILES) build/windows-amd64/
	mkdir --parents release
	$(TAR) --directory build/linux-amd64 --file \
		release/$(NAME)-$(VERSION)-linux-amd64.tar.gz .
	$(TAR) --directory build/macos-amd64 --file \
		release/$(NAME)-$(VERSION)-macos-amd64.tar.gz .
	cd build/windows-amd64 && \
		zip ../../release/$(NAME)-$(VERSION)-windows-amd64.zip *

clean:
	rm --recursive build/linux-amd64/* build/macos-amd64/* build/windows-amd64/*
