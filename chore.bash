#!/bin/bash

find_container_engine() {
  local engine

  # Check if MARKISA_CONTAINER_ENGINE is set
  if [[ -n "${MARKISA_CONTAINER_ENGINE}" ]]; then
    engine="${MARKISA_CONTAINER_ENGINE}"
  fi

  # If not set, try to find docker
  if [[ -z "${engine}" ]]; then
    engine=$(command -v docker 2>/dev/null)
  fi

  # If docker is not found, try to find podman
  if [[ -z "${engine}" ]]; then
    engine=$(command -v podman 2>/dev/null)
  fi

  # If neither is found, print an error message and exit
  if [[ -z "${engine}" ]]; then
    echo "Error: no container engine found in path" >&2
    exit 1
  fi

  echo "${engine}"
}

main() {
  CONTAINER_ENGINE=$(find_container_engine)
  echo "Using ${CONTAINER_ENGINE}"

  build() {
    clean
    mkdir -p "build"
    go build -o build/apparmor cmd/apparmor/main.go
    build/apparmor
    "${CONTAINER_ENGINE}" build . -t "quay.io/iklabib/markisa" -f "containerfiles/containerfile"
  }

  setup() {
    if [ "$UID" -eq 0 ]; then
      "$@"
    else
      sudo apparmor_parser -Kr markisa.cfg
    fi 
  }

  run() {
    # TODO: do something about selinux
    docker run --rm -it -e BASE_URL=:8080 -p 8080:8080 \
           --cap-add=sys_admin \
           --cap-add=sys_chroot \
           --cap-add=sys_resource \
           --security-opt apparmor=markisa \
           --security-opt label=disable quay.io/iklabib/markisa
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