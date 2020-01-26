FROM golang:1.13.5-alpine3.11

RUN mkdir /go
WORKDIR /go

RUN apk --update add tzdata git alpine-sdk && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

RUN go get -u github.com/oxequa/realize