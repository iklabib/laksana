FROM docker.io/library/golang:alpine3.18 as build

ENV GOPATH /
WORKDIR /build
COPY instance/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module markisa\ngo 1.21\n" > go.mod
RUN go build -o instance

FROM docker.io/library/busybox:stable
WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]
