FROM ubuntu:24.04
copy . /opt/jinghong

RUN  apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*


ENTRYPOINT ["/opt/jinghong/bin/jinghong"]
CMD ["-c", "/opt/jinghong/configs/default.yaml"]
