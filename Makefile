# run from repository root

# Example:
#   make build
#   make clean

TARGET=./bin/go-skel
MAINFILE=main.go

.PHONY: build
build:
	rm -f ${TARGET}
	MAINFILE=${MAINFILE} ./build go-skel ${TARGET}

clean:
	rm -rf ./bin *.rpm
	./release/rpm/make clean
	cd release/rpm && ./make clean

rpm:
	./release/rpm/make
