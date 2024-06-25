#!/bin/bash

SCRIPT_DIR=$(dirname "$0")
source "$SCRIPT_DIR/goenv.bash"

main() {
	build() {
		go build -o /dev/null
	}

	buildTest() {
		go test -c -o /dev/null
	}

	run() {
		export GOMAXPROCS=1
		go test --short --json
		unset GOMAXPROCS
	}

	if [[ "$1" == "build" ]]; then
		build
	elif [[ "$1" == "build-test" ]]; then
		buildTest
	elif [[ "$1" == "run" ]]; then
		run
	else
		echo "Usage: $0 {build|build-test|run}"
		exit 1
	fi
}

main "$@"
