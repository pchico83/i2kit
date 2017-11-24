FROM golang:1.9.0-alpine3.6 as builder

RUN apk update && \
    apk add --no-cache \
        git

RUN go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit

ENV SRC_DIR=/go/src/github.com/pchico83/i2kit
WORKDIR $SRC_DIR
ADD . $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/i2kit

FROM alpine:3.6
RUN apk update && \
    apk add --no-cache \
        git \
        docker \
        less \
        groff \
        py-pip \
        qemu-img \
        qemu-system-x86_64

RUN pip install awscli
COPY --from=builder /go/bin/linuxkit /usr/local/bin/linuxkit
COPY --from=builder /usr/local/bin/i2kit /usr/local/bin/i2kit
ADD . /root/i2kit

WORKDIR /root/i2kit
ENTRYPOINT ["/usr/local/bin/i2kit"]
