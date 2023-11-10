FROM docker.io/library/golang:1.22.0-alpine3.19 as build

ENV GOPATH /
WORKDIR /build
COPY cmd/builder/csharp/main.go  ./
COPY model/ ./model
COPY util/ ./util/
RUN printf "module gitlab.com/iklabib/markisa\ngo 1.22.0\n" > go.mod
RUN go mod tidy
RUN go build -o instance

FROM mcr.microsoft.com/dotnet/sdk:8.0-alpine3.19

RUN apk add --no-cache clang build-base zlib-dev

WORKDIR /tmp/csharp/
RUN dotnet new console -f net8.0 
COPY assets/csproj.txt csharp.csproj
RUN dotnet restore

WORKDIR /box
COPY --from=build /build/instance .

ENTRYPOINT [ "/box/instance" ]