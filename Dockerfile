FROM golang:1.7.3 as builder
WORKDIR /go/src/github.com/pchico83/i2kit
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o i2kit .
RUN go test `go list ./... | grep -v vendor`

FROM alpine:latest
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/pchico83/i2kit /usr/local/bin/i2kit
ENTRYPOINT ["i2kit"]
