.PHONY: all build run clean install uninstall

all: build

build:
	go build -o tetris ./cmd/tetris

run: build
	./tetris

clean:
	go clean
	rm -f tetris

install:
	./scripts/install.sh

uninstall:
	./scripts/uninstall.sh