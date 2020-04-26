FROM debian:jessie-slim as base
RUN apt-get update && apt-get install -y ca-certificates wget tzdata

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
    wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM scratch

ARG APP_PORT

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
ADD bidder-svc /

EXPOSE $APP_PORT
CMD ["/bidder-svc"]
