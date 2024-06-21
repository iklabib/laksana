#!/bin/bash

main() {
	build() {
		clean
		go mod download
		go mod vendor
		go run cmd/apparmor/main.go
		podman build . -t "quay.io/iklabib/laksana" -f "containerfiles/containerfile"
	}

	setup() {
		if [ "$UID" -eq 0 ]; then
			"$@"
		else
			sudo cp markisa.cfg /etc/apparmor.d/markisa
			sudo aa-enforce /etc/apparmor.d/markisa
			sudo apparmor_parser -Kr /etc/apparmor.d/markisa
		fi
	}

	run() {
		podman run --rm -it -p 8000:8000 \
			--cap-add sys_admin \
			--cap-add sys_resource \
			--security-opt seccomp=profiles/seccomp/markisa.json quay.io/iklabib/markisa
	}

	clean() {
		rm -rf build
	}

	if [[ "$1" == "build" ]]; then
		build
	elif [[ "$1" == "setup" ]]; then
		setup
	elif [[ "$1" == "run" ]]; then
		run
	else
		echo "Usage: $0 {build|setup|run}"
		exit 1
	fi
}

main "$@"
