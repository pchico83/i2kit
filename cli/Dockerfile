FROM golang:1.9.0-alpine3.6 as builder
ENV SRC_DIR=/go/src/github.com/pchico83/i2kit/cli/
WORKDIR $SRC_DIR
ADD . $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/i2kit

FROM alpine:3.6
RUN apk update && \
    apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/i2kit /usr/local/bin/i2kit

WORKDIR /root/i2kit
ENTRYPOINT ["/usr/local/bin/i2kit"]
