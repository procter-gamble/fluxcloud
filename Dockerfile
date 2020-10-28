FROM golang:alpine as BUILDER
RUN apk update && apk add ca-certificates
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o fluxcloud ./cmd/

FROM gcr.io/distroless/static@sha256:0bc92a5c6c154ba95a7ec2439214c640ba8574bc28b23df182ae90a8cd11182c
COPY --from=BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=BUILDER /build/fluxcloud /fluxcloud
EXPOSE 3031
ENTRYPOINT [ "/fluxcloud" ]