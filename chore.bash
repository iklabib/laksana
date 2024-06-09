#!/bin/bash

main() {
  build() {
    clean
    mkdir -p "build"
    go build -o build/apparmor cmd/apparmor/main.go
    build/apparmor
    docker build . -t "quay.io/iklabib/markisa" -f "containerfiles/containerfile"
  }

  setup() {
    if [ "$UID" -eq 0 ]; then
      "$@"
    else
      sudo cp markisa.cfg /etc/apparmor.d/markisa
      sudo aa-enforce /etc/apparmor.d/markisa
      sudo apparmor_parser -Kr /etc/apparmor.d/markisa
      sudo apparmor_parser -Kr markisa.cfg
    fi 
  }

  run() {
    docker run --rm -it -e BASE_URL=:8080 -p $BASE_URL:8080 \
           --cap-add sys_admin \
           --cap-add sys_resource \
           --security-opt apparmor=markisa quay.io/iklabib/markisa
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