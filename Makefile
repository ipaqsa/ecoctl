clean:
	rm -fr bin

build: clean
	go build -o bin/ecoctl cmd/main.go