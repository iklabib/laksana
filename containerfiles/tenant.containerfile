FROM docker.io/library/golang:alpine3.20 as build
ENV GOPATH /go
WORKDIR /build
COPY cmd/commander/main.go .
COPY . .
RUN go mod tidy
RUN go build -o commander

FROM docker.io/library/alpine:3.20

RUN apk add --no-cache musl musl-dev clang bubblewrap git python3 py3-pip py3-pytest
RUN apk add --no-cache neovim bash 
RUN wget -qO- https://go.dev/dl/go1.22.3.linux-amd64.tar.gz | tar -xz -C /usr/local && ln -s /usr/local/go/bin/go /bin/go 

ENV USER user
ENV HOME /home/user
RUN adduser -D --home $HOME $USER && mkdir -p ${HOME}/.local/bin
COPY --from=build /build/commander ${HOME}/.local/bin 
COPY configs/minijail.cfg ${HOME}/.local/bin 
RUN chown -R user:user ${HOME}

ENTRYPOINT [ "bash" ]