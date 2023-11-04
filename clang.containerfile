FROM docker.io/library/golang:alpine3.19 as build

ENV GOPATH /
WORKDIR /build
COPY cmd/builder/c/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module gitlab.com/iklabib/markisa\ngo 1.21.6\n" > go.mod
RUN go mod tidy
RUN go build -o instance

FROM docker.io/library/alpine:3.19

RUN apk add --no-cache clang

WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]