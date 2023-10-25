FROM docker.io/library/golang AS build

WORKDIR /box

ENV GOPATH /
ENV PATH $GOPATH/bin:$PATH

COPY . .
RUN go build

FROM debian:bookworm-slim
ARG compiler_file

RUN apt update && apt install unzip

WORKDIR /tmp
COPY artifacts/$compiler_file .
RUN unzip $compiler_file -d / && rm $compiler_file 
ENV PATH="/kotlinc:${PATH}"

WORKDIR /box
COPY --from=build instance .
ENTRYPOINT ["instance"]
