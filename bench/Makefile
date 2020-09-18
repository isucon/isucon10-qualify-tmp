MAKEFILE_DIR:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))

all: run/build

run/build: **/*.go
	go build -o bench ${MAKEFILE_DIR}cmd/bench/bench.go
