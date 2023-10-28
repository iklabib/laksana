FROM docker.io/library/golang:alpine3.18 as build

ENV GOPATH /
WORKDIR /build
COPY cmd/instance/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module gitlab.com/iklabib/markisa\ngo 1.21.5\n" > go.mod
RUN go mod tidy
RUN go build -o instance

FROM docker.io/library/busybox:stable
WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]
