FROM ubuntu:18.04

RUN apt-get -y update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y \
    autoconf \
    bison \
    flex \
    gcc \
    g++ \
    git \
    libprotobuf-dev \
    libnl-route-3-dev \
    libtool \
    make \
    pkg-config \
    protobuf-compiler \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/google/nsjail.git

RUN cd /nsjail && make && mv /nsjail/nsjail /bin && rm -rf -- /nsjail
ENTRYPOINT [ "bash" ]