build: clean
	mkdir "build"
	go build -o build cmd/markisa/main.go
	podman build . -t "markisa:clang" -f "containerfiles/builder/clang.containerfile"
	podman build . -t "markisa:csharp" -f "containerfiles/builder/csharp.containerfile"
	podman build . -t "markisa:common" -f "containerfiles/runner/common.containerfile"

clean:
	rm -r "build"