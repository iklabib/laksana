FROM docker.io/library/golang:1.21.6-alpine3.19 as build

ENV GOPATH /
WORKDIR /build
COPY cmd/instance/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module gitlab.com/iklabib/markisa\ngo 1.21.6\n" > go.mod
RUN go mod tidy
RUN go build -o instance

FROM mcr.microsoft.com/dotnet/runtime:8.0-alpine

# RUN apk add --no-cache libstdc++ libgcc musl

WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]
