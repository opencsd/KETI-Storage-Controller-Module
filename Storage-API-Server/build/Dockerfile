FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/storage-data-apiserver \
    USER_UID=root \
    USER_NAME=storage-data-apiserver

COPY _output/bin/storage-data-apiserver ${OPERATOR}

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
