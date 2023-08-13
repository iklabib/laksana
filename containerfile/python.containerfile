FROM docker.io/library/golang AS build

WORKDIR /instance

ENV GOPATH /
ENV PATH $GOPATH/bin:$PATH

COPY . .
RUN go build

FROM debian:bookworm-slim

RUN apt update && apt install -y python3 

WORKDIR /sandbox
COPY --from=build /instance/instance instance

ENTRYPOINT [ "/sandbox/instance" ]
