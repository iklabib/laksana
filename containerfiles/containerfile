FROM docker.io/library/golang:alpine3.20 as build
ENV GOPATH /go
WORKDIR /build
COPY cmd/commander/main.go .
COPY . .
RUN go mod tidy
RUN go build -o commander

RUN apk add --no-cache gcc bash libcap-dev alpine-sdk
RUN git clone --depth 1 --branch linux-v2024.05.22 https://chromium.googlesource.com/chromiumos/platform/minijail && cd minijail && make 

FROM docker.io/library/alpine:3.20

COPY --from=build /build/minijail/minijail0 /usr/bin/minijail0
COPY --from=build /build/minijail/libminijail.so /usr/lib/libminijail.so
COPY --from=build /build/minijail/libminijailpreload.so /lib/libminijailpreload.so
COPY --from=build /build/minijail/libminijail.h /usr/include/libminijail.h
COPY --from=build /build/minijail/minijail0.1 /usr/share/man/man1/minijail0.1
COPY --from=build /build/minijail/minijail0.5 /usr/share/man/man5/minijail0.5

RUN chmod 755 /usr/bin/minijail0 \
    && chmod 755 /usr/lib/libminijail.so \
    && chmod 755 /lib/libminijailpreload.so \
    && chmod 644 /usr/include/libminijail.h \
    && chmod 644 /usr/share/man/man1/minijail0.1 \
    && chmod 644 /usr/share/man/man5/minijail0.5

RUN apk add --no-cache musl musl-dev clang bubblewrap git python3 py3-pip py3-pytest libcap neovim bash 
RUN wget -qO- https://go.dev/dl/go1.22.3.linux-amd64.tar.gz | tar -xz -C /usr/local && ln -s /usr/local/go/bin/go /bin/go 

ENV USER user
ENV HOME /home/user
RUN adduser -D --home $HOME $USER && mkdir -p ${HOME}/.local/bin
COPY --from=build /build/commander ${HOME}/.local/bin 
COPY configs/minijail.cfg ${HOME}/.local/bin/configs/minijail.cfg 
RUN chown -R user:user ${HOME}

ENTRYPOINT [ "bash" ]