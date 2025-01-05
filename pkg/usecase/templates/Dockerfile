# TODO: Add tag or commit hash to the base image
FROM ghcr.io/secmon-lab/alertchain

COPY policy /policy

WORKDIR /
EXPOSE 8080

ENV ALERTCHAIN_LOG_FORMAT=json
ENV ALERTCHAIN_LOG_LEVEL=info
ENV ALERTCHAIN_ADDR=0.0.0.0:8080
ENV ALERTCHAIN_POLICY_DIR=/policy

ENTRYPOINT ["/alertchain", "serve"]
