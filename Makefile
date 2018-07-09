GOCMD=go
GOBUILD=$(GOCMD) build
FLAG_TAG=-X github.com/alexhokl/rds-backup/cmd.tag=$$(git describe --abbrev=0 --tags)
FLAG_VERSION=-X github.com/alexhokl/rds-backup/cmd.version=$$(git rev-parse --short HEAD)
OUTPUT_LINUX=rds-backup.linux
OUTPUT_MAC=rds-backup.darwin
OUTPUT_WINDOWS=rds-backup.exe

default: build
build:
	go get -t -v ./...
	go test -v ./...
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_LINUX) -ldflags "$(FLAG_TAG) $(FLAG_VERSION)"
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_MAC) -ldflags "$(FLAG_TAG) $(FLAG_VERSION)"
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(OUTPUT_WINDOWS) -ldflags "$(FLAG_TAG) $(FLAG_VERSION)"
release: build-linux build-mac build-windows
install:
	go install
clean:
	rm $(OUTPUT_LINUX)
	rm $(OUTPUT_MAC)
	rm $(OUTPUT_WINDOWS)
cover:
	go test -coverprofile=cover.out ./... && go tool cover -html=cover.out

