FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/opencsd-storage-api-server \
    USER_UID=root \
    USER_NAME=opencsd-storage-api-server

COPY _output/bin/opencsd-storage-api-server ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
