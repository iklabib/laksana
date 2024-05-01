build: clean
	mkdir "build"
	go build -o build cmd/markisa/main.go
	podman build . -t "iklabib/markisa:clang" -f "containerfiles/builder/clang.containerfile"
	podman build . -t "iklabib/markisa:csharp" -f "containerfiles/builder/csharp.containerfile"
	podman build . -t "iklabib/markisa:common" -f "containerfiles/runner/common.containerfile"

clean:
	rm -rf "build"