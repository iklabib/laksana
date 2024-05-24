FROM docker.io/library/golang:alpine3.19 as build
ENV GOPATH /go
WORKDIR /build
COPY cmd/commander/main.go .
COPY . .
RUN go mod tidy
RUN go build -o commander

FROM docker.io/library/alpine:3.19

RUN apk add --no-cache musl musl-dev clang bubblewrap git python3 py3-pip py3-pytest
RUN apk add --no-cache neovim bash 

ENV USER user
ENV HOME /home/user
RUN adduser -D --home $HOME $USER && mkdir -p ${HOME}/.local/bin
COPY --from=build /build/commander ${HOME}/.local/bin 
RUN chown -R user:user ${HOME}

USER $USER
WORKDIR $HOME

RUN wget -qO- https://go.dev/dl/go1.22.3.linux-amd64.tar.gz | tar -xz
# RUN wget -qO- https://raw.githubusercontent.com/pyenv/pyenv-installer/master/bin/pyenv-installer | bash && .pyenv/bin/pyenv 
RUN echo 'export PATH=$HOME/.local/bin:$HOME/go/bin:$PATH' >> .bashrc && . .bashrc

ENTRYPOINT [ "bash" ]