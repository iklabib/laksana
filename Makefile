CONTAINER_ENGINE := $(shell echo $$MARKISA_CONTAINER_ENGINE)

ifeq ($(CONTAINER_ENGINE),)
    CONTAINER_ENGINE := $(shell command -v docker 2>/dev/null)
endif

ifeq ($(CONTAINER_ENGINE),)
    CONTAINER_ENGINE := $(shell command -v podman 2>/dev/null)
endif

ifeq ($(CONTAINER_ENGINE),)
    $(error no container engine found in path)
endif

.PHONY: build
build: clean
	# mkdir "build"
	# go build -o build cmd/markisa/main.go
	$(CONTAINER_ENGINE) build . -t "quay.io/iklabib/markisa:tenant" -f "containerfiles/tenant.containerfile"

.PHONY: clean
clean:
	rm -rf "build"
