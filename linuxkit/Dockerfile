FROM golang:1.9.0-alpine3.6

RUN apk update && \
    apk add --no-cache \
        git \
        docker \
        less \
        groff \
        py-pip \
        qemu-img \
        qemu-system-x86_64

RUN go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit
RUN pip install awscli
WORKDIR /root/i2kit
ADD push.sh /root/i2kit

CMD ["./push.sh"]
