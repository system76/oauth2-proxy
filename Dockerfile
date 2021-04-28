FROM golang:1.16-buster AS builder

ARG VERSION=v7.1.2

# Copy sources
WORKDIR $GOPATH/src/github.com/oauth2-proxy/oauth2-proxy

# Fetch original source
RUN git clone --branch ${VERSION} --depth 1 https://github.com/oauth2-proxy/oauth2-proxy .

# Fetch dependencies
RUN GO111MODULE=on go mod download

# Now pull in our code
COPY providers/ providers/

# Build binary and make sure there is at least an empty key file.
#  This is useful for GCP App Engine custom runtime builds, because
#  you cannot use multiline variables in their app.yaml, so you have to
#  build the key into the container and then tell it where it is
#  by setting OAUTH2_PROXY_JWT_KEY_FILE=/etc/ssl/private/jwt_signing_key.pem
#  in app.yaml instead.
RUN VERSION=${VERSION} make build && touch jwt_signing_key.pem

# Copy binary to alpine
FROM alpine:3.13

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/oauth2-proxy/oauth2-proxy/jwt_signing_key.pem /etc/ssl/private/jwt_signing_key.pem
COPY --from=builder /go/src/github.com/oauth2-proxy/oauth2-proxy/nsswitch.conf /etc/nsswitch.conf
COPY --from=builder /go/src/github.com/oauth2-proxy/oauth2-proxy/oauth2-proxy /bin/oauth2-proxy

USER 2000:2000

ENTRYPOINT ["/bin/oauth2-proxy"]
