# run from repository root

# Example:
#   make build
#   make clean

TARGET=./bin/go-skel

.PHONY: build
build:
	rm -f ${TARGET}
	./build go-skel ${TARGET}

clean:
	rm -rf ./bin
