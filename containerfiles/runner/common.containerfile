FROM docker.io/library/golang:1.22.0-alpine3.19 as build

ENV MARKISA_RUN_TIME=2500
ENV GOPATH /
WORKDIR /build
COPY cmd/instance/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module github.com/iklabib/markisa\ngo 1.22.0\n" > go.mod
RUN go mod tidy
RUN go build -o instance

FROM alpine:3.19
WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]
