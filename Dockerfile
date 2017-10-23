FROM golang:1.9.0-alpine3.6 as builder

RUN apk update && \
    apk add --no-cache \
        git

ENV SRC_DIR=/go/src/github.com/pchico83/i2kit
WORKDIR $SRC_DIR
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/moby/tool/src/moby
RUN go get k8s.io/api/extensions/v1beta1
RUN go get k8s.io/apimachinery/pkg/util/yaml
ADD . $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/local/bin/i2kit

FROM alpine:3.6
RUN apk update && \
    apk add --no-cache \
        less \
        groff \
        py-pip

RUN pip install awscli

WORKDIR /root/i2kit
COPY --from=builder /usr/local/bin/i2kit /usr/local/bin/i2kit
ENTRYPOINT ["/usr/local/bin/i2kit"]
