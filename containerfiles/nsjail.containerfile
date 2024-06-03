FROM docker.io/library/alpine:3.20 as build
RUN apk add --no-cache alpine-sdk make gcc bison flex protobuf git bsd-compat-headers bash linux-headers libnl3-dev protobuf-dev
RUN git clone --depth 1 --branch 3.4 https://github.com/google/nsjail.git && cd nsjail && git submodule update --init
WORKDIR /nsjail
# https://github.com/google/nsjail/issues/225
COPY protobuf.patch .
RUN git apply protobuf.patch && make && install -m755 -D nsjail /usr/bin
ENTRYPOINT [ "bash" ]