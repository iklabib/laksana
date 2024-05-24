FROM docker.io/library/alpine:3.19 as build
RUN apk add --no-cache gcc bash libcap-dev alpine-sdk
RUN git clone --depth 1 https://chromium.googlesource.com/chromiumos/platform/minijail && cd minijail && make
WORKDIR /minijail

# copy minijail0, libminijail.so, libminijailpreload.so to final image
# COPY --from=build minijail0 /bin/minijail0
# COPY --from=build libminijail.so /lib/libminijail.so
# COPY --from=build libminijailpreload.so /lib/libminijailpreload.so
# COPY --from=build libminijail.h /usr/include/libminijail.h
# COPY --from=build minijail0.1 /usr/share/man/man1/minijail0.1
# COPY --from=build minijail0.5 /usr/share/man/man5/minijail0.5

# install -m755 -D minijail0 /usr/bin/minijail0
# install -m755 -D libminijail.so /usr/lib/libminijail.so
# install -m755 -D libminijailpreload.so /usr/lib/libminijailpreload.so
# install -m644 -D libminijail.h /usr/include/libminijail.h
# install -m644 -D minijail0.1 /usr/share/man/man1/minijail0.1
# install -m644 -D minijail0.5 /usr/share/man/man5/minijail0.5

ENTRYPOINT [ "bash" ]