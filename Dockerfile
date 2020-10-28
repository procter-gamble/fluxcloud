FROM alpine:3.12.1
RUN apk update && apk add ca-certificates

FROM gcr.io/distroless/static@sha256:0bc92a5c6c154ba95a7ec2439214c640ba8574bc28b23df182ae90a8cd11182c
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY fluxcloud /fluxcloud
EXPOSE 3031
ENTRYPOINT ["/fluxcloud"]