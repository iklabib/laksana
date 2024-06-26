#!/bin/bash

main() {
	_prebuild() {
		go mod download
		go mod vendor
	}

	build() {
		_prebuild
		podman build . -t "quay.io/iklabib/laksana" -f "containerfiles/containerfile"
	}

	apparmor() {
		if [ "$UID" -eq 0 ]; then
			echo "do not run as root"
			exit
		fi

		_prebuild
		go run cmd/apparmor/main.go

		if [ "$UID" -eq 0 ]; then
			"$@"
		else
			sudo cp laksana.cfg /etc/apparmor.d/laksana
			sudo aa-enforce /etc/apparmor.d/laksana
			sudo apparmor_parser -Kr /etc/apparmor.d/laksana
		fi
	}

	run() {
		podman run --rm -it -p 31415:8000 \
			--cap-add sys_admin \
			--cap-add sys_resource \
			--security-opt seccomp=profiles/seccomp/laksana.json quay.io/iklabib/laksana
	}

	if [[ "$1" == "build" ]]; then
		build
	elif [[ "$1" == "apparmor" ]]; then
		apparmor
	elif [[ "$1" == "run" ]]; then
		run
	else
		echo "Usage: $0 {build|setup|run}"
		exit 1
	fi
}

main "$@"
