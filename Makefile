default: build
build:
	go get -t -v ./...
	go test -v ./...
build-linux:
	GOOS=linux GOARCH=amd64 go build -o rds-backup.linux -ldflags "-X github.com/alexhokl/rds-backup/cmd.tag=$$(git describe --abbrev=0 --tags) -X github.com/alexhokl/rds-backup/cmd.version=$$(git rev-parse --short HEAD)"
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o rds-backup.darwin -ldflags "-X github.com/alexhokl/rds-backup/cmd.tag=$$(git describe --abbrev=0 --tags) -X github.com/alexhokl/rds-backup/cmd.version=$$(git rev-parse --short HEAD)"
build-windows:
	GOOS=windows GOARCH=amd64 go build -o rds-backup.exe -ldflags "-X github.com/alexhokl/rds-backup/cmd.tag=$$(git describe --abbrev=0 --tags) -X github.com/alexhokl/rds-backup/cmd.version=$$(git rev-parse --short HEAD)"
release: build-linux build-mac build-windows
install:
	go install
clean:
	rm rds-backup.darwin
	rm rds-backup.linux
	rm rds-backup.exe

