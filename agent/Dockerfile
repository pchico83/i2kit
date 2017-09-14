FROM golang:1.9.0-alpine3.6 as builder
ENV SRC_DIR=/go/src/bitbucket.org/riberaproject/agent
WORKDIR $SRC_DIR
ADD . $SRC_DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o agent main.go

FROM alpine:3.6
RUN apk update && \
    apk add --no-cache \
        py-pip
RUN pip install 'docker-compose==1.17.1'

WORKDIR /root/i2kit
COPY --from=builder /go/src/bitbucket.org/riberaproject/agent/agent .
CMD ["./agent"]
