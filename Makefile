.PHONY: all build

all: build

build:
	go build -o bin/elmobot ./cmd/elmo

clean:
	rm -rf bin

docker-build:
	docker build -t elmobot:latest .

docker-run:
	docker run --rm -it elmobot:latest

format:
	go fmt ./...
