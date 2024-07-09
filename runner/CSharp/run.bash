#!/bin/bash

main() {
    build() {
        output/LittleRosie build main.dll
    }

    run() {
        output/LittleRosie run main.dll
    }
	
	if [[ "$1" == "build" ]]; then
		build
	elif [[ "$1" == "run" ]]; then
		run
	else
		echo "Usage: $0 {build|run}"
		exit 1
	fi
}

main "$@"