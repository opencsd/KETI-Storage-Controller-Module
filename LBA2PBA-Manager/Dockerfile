FROM ubuntu:18.04

ADD ./cmake/build/lba2pba-manager /usr/bin/lba2pba-manager

COPY shared_library/libstdc++.so.6 /usr/lib/x86_64-linux-gnu/

ENTRYPOINT ["/usr/bin/lba2pba-manager"]
