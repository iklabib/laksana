FROM docker.io/library/alpine:3.20 as build
RUN apk add --no-cache clang make libcap-dev
ADD https://chromium.googlesource.com/chromiumos/platform/minijail/+archive/99e8fd4bf9aaf62eab9b3cabddc2939cb3427029.tar.gz minijail.tar.gz
RUN mkdir minijail && tar -xzf minijail.tar.gz -C minijail && rm minijail.tar.gz && cd minijail && make

FROM docker.io/library/alpine:3.20
RUN apk add --no-cache libcap bash

COPY --chmod=775 --from=build /minijail/minijail0 /usr/bin/minijail0
COPY --chmod=775 --from=build /minijail/libminijail.so /usr/lib/libminijail.so
COPY --chmod=775 --from=build /minijail/libminijailpreload.so /lib/libminijailpreload.so
COPY --chmod=644 --from=build /minijail/libminijail.h /usr/include/libminijail.h
COPY --chmod=644 --from=build /minijail/minijail0.1 /minijail/minijail0.1
COPY --chmod=644 --from=build /minijail/minijail0.5 /minijail/minijail0.5

ENV USER user
ENV HOME /home/user
RUN adduser -D --home $HOME $USER 
RUN chown -R user:user ${HOME}

ENTRYPOINT [ "bash" ]