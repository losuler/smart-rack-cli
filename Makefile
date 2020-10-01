NAME=srs
VERSION=$(shell git describe)
TAR=tar --create --verbose --file

build:
	go build -o bin/$(NAME) main.go

run:
	go run main.go

build-all:
	GOOS=linux GOOARCH=amd64 go build -v -o build/linux-amd64/$(NAME) main.go
	GOOS=darwin GOARCH=amd64 go build -v -o build/macos-amd64/$(NAME) main.go
	GOOS=windows GOARCH=amd64 go build -v -o build/windows-amd64/$(NAME) main.go

release-all:
	$(TAR) $(NAME)-$(VERSION)-linux-amd64.tar.gz build/linux-amd64/$(NAME)
	$(TAR) $(NAME)-$(VERSION)-macos-amd64.tar.gz build/macos-amd64/$(NAME)
	$(TAR) $(NAME)-$(VERSION)-windows-amd64.tar.gz build/windows-amd64/$(NAME)

clean:
	rm --recursive build/