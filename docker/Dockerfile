FROM harness0.harness.io/oci/docker_artifacts/hhi/harnessos-alpine:1@sha256:290a37c85b4284d2c11fa0778e5da79643d912a1f5a57d6f701b77644503b7e7 AS alpine
RUN apk add -U --no-cache 'ca-certificates>=20261020-r0'

FROM harness0.harness.io/oci/docker_artifacts/hhi/harnessos-alpine:1@sha256:290a37c85b4284d2c11fa0778e5da79643d912a1f5a57d6f701b77644503b7e7
ENV GODEBUG netdns=go

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ADD release/linux/amd64/plugin /bin/
ENTRYPOINT ["/bin/plugin"]
