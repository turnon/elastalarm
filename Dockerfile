FROM golang:alpine as builder

WORKDIR $GOPATH/src/github.com/turnon/elastalarm
COPY . ./

RUN apk add --no-cache git \
    && go get ./... \
    && go build -o /elastalarm \
    && apk del git

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=builder /elastalarm .
RUN chmod +x /root/elastalarm

ENTRYPOINT ["/root/elastalarm", "-configs", "/configs"]