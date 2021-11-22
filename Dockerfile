FROM quay.io/sonraisecurity/alpine:3.14-2 AS production
USER root
COPY ./build/asearch /sonrai/bin/service
ENTRYPOINT ["service"]