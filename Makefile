.PHONY: build clean

build:
	go build -o bin/manbo *.go

clean:
	rm -rf bin/