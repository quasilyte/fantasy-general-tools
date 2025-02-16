.PHONY: build-windows build-linux

build-windows:
	mkdir -p bin/
	GOOS=windows go build -o bin/fgtool.exe -trimpath -ldflags='-s' ./cmd/fgtool

build-linux:
	mkdir -p bin/
	GOOS=linux go build -o bin/fgtool -trimpath -ldflags='-s' ./cmd/fgtool
