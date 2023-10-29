VERSION=v0.0.1

clean:
	rm -rf bin

build: clean
	go build -ldflags "-s -w -X ecoctl/pkg/version.Version=${VERSION}" -o bin/ecoctl cmd/main.go