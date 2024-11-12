FROM ubuntu:18.04

ENV OPERATOR=/usr/local/bin/storage-metric-collector \
    USER_UID=root \
    USER_NAME=storage-metric-collector

COPY bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/storage-metric-collector"]

USER ${USER_UID}
