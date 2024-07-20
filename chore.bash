#!/bin/bash

main() {
	_prebuild() {
		go mod download
		go mod vendor
	}

	build() {
		_prebuild
		${CONTAINER_ENGINE} build . -t "quay.io/iklabib/laksana" -f "containerfiles/containerfile"
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
		if [[ "$1" == "--apparmor" ]]; then
			${CONTAINER_ENGINE} run --rm -it -p 31415:8000 \
				--cap-add sys_admin \
				quay.io/iklabib/laksana
		else
			${CONTAINER_ENGINE} run --rm -it -p 31415:8000 \
				--cap-add sys_admin \
				quay.io/iklabib/laksana
		fi
	}

	if [[ -z "${LAKSANA_CE}" ]]; then
		CONTAINER_ENGINE="docker"
	else
		CONTAINER_ENGINE="${LAKSANA_CE}"
	fi

	if [[ "$1" == "build" ]]; then
		build
	elif [[ "$1" == "apparmor" ]]; then
		apparmor
	elif [[ "$1" == "run" ]]; then
		run
	else
		echo "Usage: $0 {build|apparmor|run [--apparmor]}"
		exit 1
	fi
}

main "$@"
