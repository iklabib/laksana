FROM docker.io/library/alpine:3.20 as build
RUN apk add --no-cache gcc bash libcap-dev alpine-sdk
RUN git clone --depth 1 https://chromium.googlesource.com/chromiumos/platform/minijail && cd minijail && make 

FROM docker.io/library/alpine:3.20
# copy minijail0, libminijail.so, libminijailpreload.so to final image
COPY --from=build /minijail/minijail0 /usr/bin/minijail0
COPY --from=build /minijail/libminijail.so /usr/lib/libminijail.so
COPY --from=build /minijail/libminijailpreload.so /usr/lib/libminijailpreload.so
COPY --from=build /minijail/libminijail.h /usr/include/libminijail.h
COPY --from=build /minijail/minijail0.1 /usr/share/man/man1/minijail0.1
COPY --from=build /minijail/minijail0.5 /usr/share/man/man5/minijail0.5

RUN chmod 755 /usr/bin/minijail0 \
    && chmod 755 /usr/lib/libminijail.so \
    && chmod 755 /usr/lib/libminijailpreload.so \
    && chmod 644 /usr/include/libminijail.h \
    && chmod 644 /usr/share/man/man1/minijail0.1 \
    && chmod 644 /usr/share/man/man5/minijail0.5

RUN apk add --no-cache libcap bash

ENTRYPOINT [ "bash" ]